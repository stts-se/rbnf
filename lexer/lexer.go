package lexer

import (
	"fmt"
	//"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Result []Item

func (r Result) String() string {
	res := []string{}
	for _, i := range r {
		res = append(res, i.String())
	}
	return strings.Join(res, ", ")
}

func (r Result) Errors() []Item {
	var res []Item
	for _, i := range r {
		if i.Typ == ItemError {
			res = append(res, i)
		}
	}
	return res
}

type Lexer struct {
	input  string // the string being scanned
	start  int    // start position of this Item
	pos    int    // current position in the input
	width  int    // width of last rune read
	Result Result // slice of scanned Items
	state  stateFn
}

type Item struct {
	Typ ItemType
	Val string
}

func (i Item) String() string {
	switch i.Typ {
	case ItemError:
		return i.Val
	}
	// if len(i.val) > 10 {
	// 	return fmt.Sprintf("%s{%.10v...}", i.typ, i.val)
	// }
	return fmt.Sprintf("%s{%v}", i.Typ, i.Val)
}

//go:generate stringer -type=ItemType

type ItemType int

const (
	ItemError ItemType = iota
	ItemLeftBracket
	ItemLeftSub
	ItemRightSub
	ItemRightBracket
	ItemSpellout
	ItemLeftDelim
	ItemRightDelim
	ItemVariable
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
	delimChars = " -\u00AD" // \u00AD = soft hyphen
	aToZ       = "abcdefghijklmnopqrstuvwxyz"
	//leftSubChars  = "←%[]" + aToZ + delimChars
	//rightSubChars = "→%[]" + aToZ + delimChars
	//spelloutChars = aToZ
	ruleNameChars = aToZ + "-"
	x             = 'x'
)

func (t ItemType) String() string {
	switch t {
	case ItemError:
		return "error"
	case ItemLeftSub:
		return "leftsub"
	case ItemRightSub:
		return "rightsub"
	case ItemLeftBracket:
		return "leftbracket"
	case ItemRightBracket:
		return "rightbracket"
	case ItemSpellout:
		return "spellout"
	case ItemLeftDelim:
		return "leftdelim"
	case ItemRightDelim:
		return "rightdelim"
	case ItemVariable:
		return "variable"
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
		i := Item{
			ItemError,
			fmt.Sprintf(format, args...),
		}
		l.Result = append(l.Result, i)
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

func spelloutFn(l *Lexer) stateFn {
	var plainSpelloutFunc = func() bool {
		var acceptFunc = func(r rune) bool {
			return unicode.IsLetter(r) || r == '\''
		}
		if l.acceptRunFunc(acceptFunc) > 0 {
			l.emit(ItemSpellout)
			return true
		}
		return false
	}
	var ruleRefFunc = func() bool {
		r, r2, r3 := l.peek3()
		if strings.IndexRune(delimChars, r) >= 0 && r2 == '=' && r3 == '%' {
			l.acceptRunString(delimChars)
			l.emit(ItemSpellout)
			r, r2, r3 = l.peek3()
		}

		if r == '=' && r2 == '%' && strings.IndexRune(ruleNameChars, r3) >= 0 {
			l.next()
			l.next()
			l.acceptRunString(ruleNameChars)
			if l.next() == '=' {
				l.emit(ItemSpellout)
				return true
			}
		}
		return false
	}
	for {
		if !(plainSpelloutFunc() || ruleRefFunc()) {
			break
		}

		r, r2, _ := l.peek3()
		if r == rightArr || r == leftBracket || (l.acceptPeekString(delimChars) && r2 == rightArr) {
			return rightSubFn
		}
		if l.acceptRunString(delimChars) > 0 {
			l.emit(ItemSpellout)
		}
	}
	switch l.peek() {
	case endTag:
		return nil
	case eof:
		return prematureEndOfInput
	default:
		return l.errorf("unknown input at expected %s: '%s'", ItemSpellout, l.currentToEnd())
	}
	return nil
}

func leftSubFn(l *Lexer) stateFn {
	closingTag := leftArr

	// opening tags
	for {
		r := l.peek()
		if r == leftBracket {
			l.next()
			l.emit(ItemLeftBracket)
			closingTag = rightBracket
			break
		} else if r == leftArr {
			l.next()
			break
		} else {
			return l.errorf("unknown opening input at expected %s: '%s'", ItemLeftSub, l.currentToEnd())
		}
	}

	for {
		r := l.peek()
		if r == closingTag {
			if r == leftArr {
				l.next()
				l.emit(ItemLeftSub)
				if l.acceptRunString(delimChars) > 0 {
					l.emit(ItemLeftDelim)
				}
				return spelloutFn
			} else if r == rightBracket {
				if len(l.current()) > 0 {
					l.emit(ItemLeftSub)
				}
				l.next()
				l.emit(ItemRightBracket)
				return spelloutFn
			}
		} else if r == leftArr {
			l.next()
		} else if l.acceptPeekString(delimChars) {
			l.emit(ItemLeftSub)
			if l.acceptRunString(delimChars) > 0 {
				l.emit(ItemLeftDelim)
			}
		} else if r == '%' {
			l.next()
			l.acceptRunString(ruleNameChars)
		} else {
			return l.errorf("unknown input at expected %s: '%s'", ItemLeftSub, l.currentToEnd())
		}
	}
	panic("not reached")
}

func rightSubFn(l *Lexer) stateFn {
	closingTag := rightArr

	// opening tags
	for {
		r := l.peek()
		if r == leftBracket {
			l.next()
			l.emit(ItemLeftBracket)
			closingTag = rightBracket
			break
		} else if r == rightArr {
			l.next()
			break
		} else if l.acceptPeekString(delimChars) {
			l.next()
			l.emit(ItemRightDelim)
		} else {
			return l.errorf("unknown opening input at expected %s: '%s'", ItemRightSub, l.currentToEnd())
		}
	}

	for {
		r := l.peek()
		if r == closingTag {
			if r == rightArr {
				l.next()
				l.emit(ItemRightSub)
				return endFn
			} else if r == rightBracket {
				l.emit(ItemRightSub)
				l.next()
				l.emit(ItemRightBracket)
				return endFn
			}
		} else if r == rightArr {
			l.next()
		} else if l.acceptPeekString(delimChars) {
			//l.next()
			if l.acceptRunString(delimChars) > 0 {
				l.emit(ItemRightDelim)
			}
		} else if r == '%' {
			l.next()
			l.acceptRunString(ruleNameChars)
		} else {
			return l.errorf("unknown input at expected %s: '%s'", ItemLeftSub, l.currentToEnd())
		}
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

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		state:  initialState,
		Result: Result{},
	}
	return l
}

func (l *Lexer) current() string {
	return l.input[l.start:l.pos]
}
func (l *Lexer) currentToEnd() string {
	return l.input[l.pos:]
}
func (l *Lexer) emit(t ItemType) {
	i := Item{t, l.current()}
	//fmt.Printf("emit: %s\n", i)
	l.Result = append(l.Result, i)
	l.start = l.pos
}

func (l *Lexer) debug(msg string) {
	s := l.current()
	fmt.Printf("lexer debug %s: '%s'\n", msg, s)
	//fmt.Printf("lexer debug %s: '%#v'\n", msg, l)
}

func (l *Lexer) Run() error {
	for state := initialState; state != nil; {
		state = state(l)
	}
	for _, i := range l.Result {
		if i.Typ == ItemError {
			return fmt.Errorf("%v", i.Val)
		}
	}
	return nil
}

// func main() {
// 	for _, s := range os.Args[1:] {
// 		fmt.Printf("input: '%s'\n", s)
// 		l := Lex(s)
// 		l.Run()
// 		for _, i := range l.Result {
// 			fmt.Printf("item: %s\n", i)
// 		}
// 		fmt.Println()
// 	}
// }
