package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type result []item

func (r result) String() string {
	res := []string{}
	for _, i := range r {
		res = append(res, i.String())
	}
	return strings.Join(res, ", ")
}

type lexer struct {
	input  string // the string being scanned
	start  int    // start position of this item
	pos    int    // current position in the input
	width  int    // width of last rune read
	result result // slice of scanned items
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
	// if len(i.val) > 10 {
	// 	return fmt.Sprintf("%s{%.10v...}", i.typ, i.val)
	// }
	return fmt.Sprintf("%s{%v}", i.typ, i.val)
}

//go:generate stringer -type=itemType

type itemType int

const (
	itemError itemType = iota
	itemLeftBracket
	itemLeftSub
	itemRightSub
	itemRightBracket
	itemSpellout
	itemDelim
	itemVariable
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
	delimChars = " -"
	aToZ       = "abcdefghijklmnopqrstuvwxyz"
	//leftSubChars  = "←%[]" + aToZ + delimChars
	//rightSubChars = "→%[]" + aToZ + delimChars
	spelloutChars = aToZ
	ruleNameChars = aToZ + "-"
	x             = 'x'
)

func (t itemType) String() string {
	switch t {
	case itemError:
		return "error"
	case itemLeftSub:
		return "leftsub"
	case itemRightSub:
		return "rightsub"
	case itemLeftBracket:
		return "leftbracket"
	case itemRightBracket:
		return "rightbracket"
	case itemSpellout:
		return "spellout"
	case itemDelim:
		return "delim"
	case itemVariable:
		return "variable"
	default:
		panic(fmt.Sprintf("undefined string output for %d", t))
	}
}

type stateFn func(*lexer) stateFn

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	return func(lx *lexer) stateFn {
		i := item{
			itemError,
			fmt.Sprintf(format, args...),
		}
		l.result = append(l.result, i)
		return nil
	}
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) back() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.back()
	return r
}

func (l *lexer) acceptRun(valid string) int {
	n := 0
	for strings.IndexRune(valid, l.next()) >= 0 {
		n++
	}
	l.back()
	return n
}

func (l *lexer) acceptPeek(valid string) bool {
	return strings.IndexRune(valid, l.peek()) >= 0
}

func spelloutFn(l *lexer) stateFn {
	if l.acceptRun(spelloutChars) > 0 {
		l.emit(itemSpellout)
	}
	r := l.peek()
	if r == rightArr || r == leftBracket || l.acceptPeek(delimChars) {
		return rightSubFn
	}
	switch l.peek() {
	case endTag:
		return nil
	case eof:
		return prematureEndOfInput
	default:
		return l.errorf("unknown input at expected %s: '%s'", itemSpellout, l.currentToEnd())
	}
	return nil
}

func leftSubFn(l *lexer) stateFn {
	closingTag := leftArr

	// opening tags
	for {
		r := l.peek()
		if r == leftBracket {
			l.next()
			l.emit(itemLeftBracket)
			closingTag = rightBracket
			break
		} else if r == leftArr {
			l.next()
			break
		} else {
			return l.errorf("unknown opening input at expected %s: '%s'", itemLeftSub, l.currentToEnd())
		}
	}

	for {
		r := l.peek()
		if r == closingTag {
			if r == leftArr {
				l.next()
				l.acceptRun(delimChars)
				l.emit(itemLeftSub)
				return spelloutFn
			} else if r == rightBracket {
				l.emit(itemLeftSub)
				l.next()
				l.emit(itemRightBracket)
				return spelloutFn
			}
		} else if r == leftArr {
			l.next()
		} else if l.acceptPeek(delimChars) {
			l.next()
		} else if r == '%' {
			l.next()
			l.acceptRun(ruleNameChars)
		} else {
			return l.errorf("unknown input at expected %s: '%s'", itemLeftSub, l.currentToEnd())
		}
	}
	panic("not reached")
}

func rightSubFn(l *lexer) stateFn {
	closingTag := rightArr

	// opening tags
	for {
		r := l.peek()
		if r == leftBracket {
			l.next()
			l.emit(itemLeftBracket)
			closingTag = rightBracket
			break
		} else if r == rightArr {
			l.next()
			break
		} else if l.acceptPeek(delimChars) {
			l.next()
		} else {
			return l.errorf("unknown opening input at expected %s: '%s'", itemRightSub, l.currentToEnd())
		}
	}

	for {
		r := l.peek()
		if r == closingTag {
			if r == rightArr {
				l.next()
				l.emit(itemRightSub)
				return endFn
			} else if r == rightBracket {
				l.emit(itemRightSub)
				l.next()
				l.emit(itemRightBracket)
				return endFn
			}
		} else if r == rightArr {
			l.next()
		} else if l.acceptPeek(delimChars) {
			l.next()
		} else if r == '%' {
			l.next()
			l.acceptRun(ruleNameChars)
		} else {
			return l.errorf("unknown input at expected %s: '%s'", itemLeftSub, l.currentToEnd())
		}
	}
	panic("not reached")
}

func prematureEndOfInput(l *lexer) stateFn {
	return l.errorf("premature end of input")
}

func endFn(l *lexer) stateFn {
	switch r := l.peek(); {
	case r == endTag:
		l.next()
		return nil
	default:
		return l.errorf("unknown input at expected %s: '%s'", "end", l.currentToEnd())
	}
}

func initialState(l *lexer) stateFn {
	switch r := l.peek(); {
	case r == leftArr || r == leftBracket:
		return leftSubFn
	case r == endTag:
		l.next()
		return nil
	case r == eof:
		l.next()
		return prematureEndOfInput
	default:
		return spelloutFn
	}
	return nil
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		state:  initialState,
		result: result{},
	}
	return l
}

func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}
func (l *lexer) currentToEnd() string {
	return l.input[l.pos:]
}
func (l *lexer) emit(t itemType) {
	i := item{t, l.current()}
	//fmt.Printf("emit: %s\n", i)
	l.result = append(l.result, i)
	l.start = l.pos
}

func (l *lexer) debug(msg string) {
	s := l.current()
	fmt.Printf("lexer debug %s: '%s'\n", msg, s)
	//fmt.Printf("lexer debug %s: '%#v'\n", msg, l)
}

func (l *lexer) run() {
	for state := initialState; state != nil; {
		state = state(l)
	}
}

func main() {
	for _, s := range os.Args[1:] {
		fmt.Printf("input: '%s'\n", s)
		l := lex(s)
		l.run()
		for _, i := range l.result {
			fmt.Printf("item: %s\n", i)
		}
		fmt.Println()
	}
}
