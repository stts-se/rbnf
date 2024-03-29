/** Package lexer contains a parser for the ICU project rbnf rule format. The code in the file lexer.go is based on the template lexer of the standard Go distribution: https://golang.org/src/text/template/parse/lex.go

A presentation of the original code:
https://talks.golang.org/2011/lex.slide
https://www.youtube.com/watch?v=HxaD_trXwRE

The original code is published under a BSD license: https://golang.org/LICENSE
*/

package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type result []item

func (r result) String() string {
	res := []string{}
	for _, i := range r {
		res = append(res, i.String())
	}
	return strings.Join(res, ", ")
}

func (r result) Errors() []item {
	var res []item
	for _, i := range r {
		if i.typ == itemError {
			res = append(res, i)
		}
	}
	return res
}

// Lexer is a struct for used to parse an input rbnf string into a slice of items
type Lexer struct {
	input  string // the string being scanned
	start  int    // start position of this Item
	pos    int    // current position in the input
	width  int    // width of last rune read
	result result // slice of scanned Items
	state  stateFn
}

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch i.typ {
	case itemError:
		return i.val
	}
	return fmt.Sprintf("%s{%v}", i.typ, i.val)
}

//go:generate stringer -type=ItemType (doesn't work with go mod?)

type itemType int

const (
	itemError itemType = iota
	itemSub
	itemLeftBracket
	itemRightBracket
)

const (
	// rune constants
	eof          = 25 // ASCII end of medium
	leftArr      = '←'
	rightArr     = '→'
	leftBracket  = '['
	rightBracket = ']'
	endTag       = ';'

	// string constants
	rulePointer             = "←→="
	aToZ                    = "abcdefghijklmnopqrstuvwxyz"
	delimCharsNotInRuleName = "', \u00AD−" // \u00AD = soft hyphen
	delimChars              = delimCharsNotInRuleName + "-"
	ruleNameChars           = aToZ + "-"
	x                       = 'x'
)

func (t itemType) String() string {
	switch t {
	case itemError:
		return "error"
	case itemSub:
		return "sub"
	case itemLeftBracket:
		return "leftbracket"
	case itemRightBracket:
		return "rightbracket"
	default:
		panic(fmt.Sprintf("undefined string output for %d", t))
	}
}

type stateFn func(*Lexer) stateFn

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Lexer) next2() (rune, rune) {
	return l.next(), l.next()
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	return func(lx *Lexer) stateFn {
		i := item{
			itemError,
			fmt.Sprintf(format, args...),
		}
		l.result = append(l.result, i)
		return nil
	}
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) back() {
	l.pos -= l.width
	// todo: should set new width (cf next)
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.back()
	return r
}

func (l *Lexer) peek3() (rune, rune, rune) {
	posBefore := l.pos
	var r1 = l.next()
	var w1 = l.width
	var r2 = rune(eof)
	var w2 = 0
	var r3 = rune(eof)
	var w3 = 0
	if r1 != eof {
		r2 = l.next()
		w2 = l.width
		if r2 != eof {
			r3 = l.next()
			w3 = l.width
		}
	}
	l.pos -= (w1 + w2 + w3)
	posAfter := l.pos
	if posAfter != posBefore {
		panic(fmt.Sprintf("??? <%s> <%s> %d %d %#v", string(r1), string(r2), posBefore, posAfter, l))
	}
	return r1, r2, r3
}

func (l *Lexer) acceptRunString(valid string) int {
	n := 0
	for strings.IndexRune(valid, l.next()) >= 0 {
		n++
	}
	l.back()
	return n
}

func (l *Lexer) acceptRunFunc(f func(r rune) bool) int {
	n := 0
	for f(l.next()) {
		n++
	}
	l.back()
	return n
}

func (l *Lexer) acceptPeekString(valid string) bool {
	return strings.IndexRune(valid, l.peek()) >= 0
}

func isPlainText(r rune) bool {
	return r != rightArr && r != leftArr && r != rightBracket && r != leftBracket && r != '=' && r != ';' && r != eof && r != '$'
}

func nfc(s string) string {
	normed, _, _ := transform.String(norm.NFC, s)
	return normed
}

func subFn(l *Lexer) stateFn {
	var closingFunc func(r rune) (bool, bool) // first bool: we're at a closing tag; 2nd bool: include current rune before emitting

	// opening tags
	for {
		r := l.peek()
		//fmt.Println("opening with", string(r))
		if r == eof {
			l.next()
			return prematureEndOfInput
			break
		} else if r == endTag {
			l.next()
			return nil
		} else if r == leftBracket {
			l.next()
			l.emit(itemLeftBracket)
			return subFn
		} else if r == rightBracket {
			l.next()
			l.emit(itemRightBracket)
			return subFn
		} else if r == rightArr || r == leftArr || r == '=' {
			l.next()
			closingFunc = func(rx rune) (bool, bool) {
				if rx == r {
					return true, true
				}
				if strings.IndexRune(l.current(), '%') >= 0 {
					if strings.IndexRune(delimCharsNotInRuleName, rx) >= 0 {
						return true, false
					}
					return false, false
				}
				if strings.IndexRune(l.current(), '#') >= 0 {
					if strings.IndexRune("#,.0123456789", rx) < 0 {
						return true, true
					}
					return false, false
				}
				if strings.IndexRune(delimChars, rx) >= 0 {
					return true, false
				}
				if strings.IndexRune("$", rx) >= 0 {
					return true, false
				}
				return false, false
			}
			break
		} else if r == '$' {
			closingFunc = func(rx rune) (bool, bool) {
				if rx == r {
					return true, true
				}
				return false, false
			}
			l.next()
			break
		} else if isPlainText(r) {
			closingFunc = func(rx rune) (bool, bool) {
				return !(isPlainText(rx)), false
			}
			l.next()
			break
		} else {
			return l.errorf("unknown opening input at expected %v: '%v'", itemSub, l.currentToEnd())
		}
	}

	for {
		r := l.peek()
		//fmt.Printf("input: '%s'\n", l.input)
		//fmt.Printf("acc result: '%s'\n", l.Result())
		//fmt.Println("inside tags", r, string(r))
		if doClose, includeClosingRune := closingFunc(r); doClose {
			//fmt.Printf("%v '%v' | doClose %v, includeClosingRune %v\n", r, string(r), doClose, includeClosingRune)
			if includeClosingRune {
				l.next()
			}
			l.emit(itemSub)
			if r == eof {
				l.next()
				return prematureEndOfInput
			}
			if r == endTag {
				l.next()
				return nil
			}
			return subFn

		}
		if r == eof {
			l.next()
			return prematureEndOfInput
		}
		l.next()
		//return l.errorf("unknown input at expected %s: '%s'", ItemSub, l.currentToEnd())
	}
	panic("not reached")
}

func prematureEndOfInput(l *Lexer) stateFn {
	return l.errorf("premature end of input")
}

func endFn(l *Lexer) stateFn {
	switch r := l.peek(); {
	case r == endTag:
		l.next()
		return nil
	default:
		return l.errorf("unknown input at expected %s: '%s'", "end", l.currentToEnd())
	}
}

func initialState(l *Lexer) stateFn {
	switch r := l.peek(); {
	case r == endTag:
		l.next()
		return nil
	case r == eof:
		l.next()
		return prematureEndOfInput
	default:
		return subFn
	}
	return nil
}

// Lex creates a new Lexer for the input string
func Lex(input string) *Lexer {
	input = nfc(input)
	l := &Lexer{
		input:  input,
		state:  initialState,
		result: result{},
	}
	return l
}

// Result creates the final result (slice of strings)
func (l *Lexer) Result() []string {
	res := []string{}
	// Post-process brackets for optional content
	openBracket := false
	for _, item := range l.result {
		if item.typ == itemLeftBracket {
			openBracket = true
		} else if item.typ == itemRightBracket {
			openBracket = false
		} else if openBracket {
			item.val = "[" + item.val + "]"
			res = append(res, item.val)
		} else {
			res = append(res, item.val)
		}
	}
	return res
}

func (l *Lexer) current() string {
	return l.input[l.start:l.pos]
}
func (l *Lexer) currentToEnd() string {
	return l.input[l.pos:]
}
func (l *Lexer) emit(t itemType) {
	i := item{t, l.current()}
	//fmt.Printf("emit: %s\n", i)
	l.result = append(l.result, i)
	l.start = l.pos
}

func (l *Lexer) debug(msg string) {
	s := l.current()
	fmt.Printf("lexer debug %s: '%s'\n", msg, s)
	//fmt.Printf("lexer debug %s: '%#v'\n", msg, l)
}

// Run is called to parse the input string into separate items
func (l *Lexer) Run() error {
	//fmt.Printf("Lexer input: '%s'\n", l.input)
	for state := initialState; state != nil; {
		state = state(l)
	}
	for _, i := range l.result {
		if i.typ == itemError {
			return fmt.Errorf("%v", i.val)
		}
	}
	return nil
}
