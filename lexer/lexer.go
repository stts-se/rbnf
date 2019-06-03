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
	ItemDelim
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
	case ItemDelim:
		return "delim"
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
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.back()
	return r
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
	if l.acceptRunFunc(unicode.IsLetter) > 0 {
		l.emit(ItemSpellout)
		// for {
		// 	r := l.peek()
		// 	if strings.IndexRune(delimChars, r) >= 0 {
		// 		l.next()
		// 		if l.acceptRunFunc(unicode.IsLetter) > 0 {
		// 			// continue looking
		// 			//l.back()
		// 			//l.emit(ItemSpellout)
		// 			//break
		// 		} else {
		// 			l.back()
		// 			break
		// 		}
		// 	} else {
		// 		l.emit(ItemSpellout)
		// 		break
		// 	}
		// }
	}

	r := l.peek()
	if r == rightArr || r == leftBracket || l.acceptPeekString(delimChars) {
		return rightSubFn
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
				l.acceptRunString(delimChars)
				l.emit(ItemLeftSub)
				return spelloutFn
			} else if r == rightBracket {
				l.emit(ItemLeftSub)
				l.next()
				l.emit(ItemRightBracket)
				return spelloutFn
			}
		} else if r == leftArr {
			l.next()
		} else if l.acceptPeekString(delimChars) {
			l.next()
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
			l.next()
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
