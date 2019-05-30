// https://www.youtube.com/watch?v=HxaD_trXwRE
// https://talks.golang.org/2011/lex.slide
package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	name  string    // used only for error reports
	input string    // the string being scanned
	start int       // start position of this item
	pos   int       // current position in the input
	width int       // width of last rune read
	items chan item // channel of scanned items
	state stateFn
}

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%s{%.10v...}", itemType2String(i.typ), i.val)
	}
	return fmt.Sprintf("%s{%v}", itemType2String(i.typ), i.val)
}

//go:generate stringer -type=itemType

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemOptionalLeftSub
	itemLeftSub
	itemRightSub
	itemSpellout
	itemPadding
	itemVariable
)

func itemType2String(t itemType) string {
	switch t {
	case itemError:
		return "error"
	case itemEOF:
		return "EOF"
	case itemLeftSub:
		return "leftsub"
	case itemOptionalLeftSub:
		return "optionalleftsub"
	case itemRightSub:
		return "rightsub"
	case itemSpellout:
		return "spellout"
	case itemPadding:
		return "padding"
	case itemVariable:
		return "variable"
	default:
		panic(fmt.Sprintf("undefined string output for %v", t))
	}
}

const (
	eof           = 25 // ASCII end of medium
	leftArr       = '←'
	rightArr      = '→'
	leftBracket   = '['
	rightBracket  = ']'
	spelloutChars = "abcdefghijklmnopqrstuvxyz"
	ruleNameChars = "abcdefghijklmnopqrstuvxyz-"
	delimChars    = "+-"
	x             = 'x'
)

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
		l.items <- i
		return nil
	}
}

func isSpace(r rune) bool {
	return r == ' '
}

func isSpelloutChar(r rune) bool {
	return strings.IndexRune(spelloutChars, r) >= 0
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

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.back()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.back()
}

func leftSubFn(l *lexer) stateFn {
	for {
		fmt.Printf("%#v\n", l)
		switch r := l.next(); {
		case r == rightArr:

		}
		if strings.HasPrefix(l.input[l.pos:], "%") {
			l.pos += len("%")
			l.accept(ruleNameChars)
		}
		l.accept(delimChars)
		if strings.HasPrefix(l.input[l.pos:], " ") {
			l.pos += len(" ")
			l.accept(" ")
		}
	}
	return l.errorf("missing end tag for: %s", item{itemLeftSub, l.input[l.start:l.pos]})
}

func startLeftSubFn(l *lexer) stateFn {
	l.pos += len(leftArr)
	return leftSubFn
}

func startLeftOptionalSubFn(l *lexer) stateFn {
	l.pos += len(leftBracket)
	l.emit(itemOptionalLeftSub)
	if strings.HasPrefix(l.input[l.pos:], leftArr) {
		l.pos += len(rightArr)
		return startLeftSubFn
	}
	return l.errorf("missing end tag for: %s", item{itemOptionalLeftSub, l.input[l.start:l.pos]})
}

func spelloutFn(l *lexer) stateFn {
	l.accept(spelloutChars)
	return l.errorf("jooo")
}

func initialState(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], leftArr) {
		return startLeftSubFn
	}
	if strings.HasPrefix(l.input[l.pos:], leftBracket) {
		return startLeftOptionalSubFn
	}

	switch r := l.next(); {
	case r == eof || r == '\n':
		break
	default:
		return spelloutFn
	}
	l.emit(itemEOF)
	return nil
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		state: initialState,
		items: make(chan item, 10), // Two should be sufficient but it doesn't work (deadlock)
	}
	return l
}

func (l *lexer) emit(t itemType) {
	i := item{t, l.input[l.start:l.pos]}
	fmt.Printf("emit: %s\n", i)
	l.items <- i
	l.start = l.pos
}

func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

func (l *lexer) run() {
	for state := initialState; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func main() {
	for _, s := range os.Args[1:] {
		l := lex("test1", s)
		l.run()
		for {
			i := l.nextItem()
			fmt.Printf("nextItem: %s\n", i)
			if i.typ == itemEOF {
				break
			}
			if i.typ == itemError {
				break
			}
		}
	}
}

// två;
// trettio[­→→];
// ←%spellout-cardinal-reale← miljoner[ →→];
// tjugo→%%ord-fem-nde→;
// minus →→;
// ←← komma →→;
// ←←­hundra[­→→];
// ←%spellout-cardinal-reale← miljon→%%ord-fem-teer→;

// made-up example
// [←%left-dummy-rule←-]komma[-→%right-dummy-rule→];
// [←%left-dummy-rule←]
// -
// komma
// -
// [→%right-dummy-rule→];
