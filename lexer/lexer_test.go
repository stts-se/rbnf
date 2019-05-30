package lexer

import (
	"fmt"
	"testing"
)

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

var fs = "For input '%s', expected '%s', got '%s'"
var fsindex = "For input '%s' item %d, expected '%s', got '%s'"

func compare(input string, exp, got Result) []error {
	res := []error{}
	if len(exp) != len(got) {
		res = append(res, fmt.Errorf(fs, input, exp, got))
		return res
	}
	for i, expi := range exp {
		goti := got[i]
		if expi != goti {
			res = append(res, fmt.Errorf(fsindex, input, i, expi, goti))
		}
	}
	return res
}

func TestBasic(t *testing.T) {
	var input string
	var exp Result
	var l *Lexer
	var prematureEOIItem = Item{itemError, "premature end of input"}

	//
	input = ""
	exp = Result{prematureEOIItem}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

	//
	input = "minus;"
	exp = Result{{itemSpellout, "minus"}}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

	//
	input = "minus"
	exp = Result{{itemSpellout, "minus"}, prematureEOIItem}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

}

func TestSub(t *testing.T) {
	var input string
	var exp Result
	var l *Lexer

	//
	input = "←← komma;"
	exp = Result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

	//
	input = "komma →→;"
	exp = Result{
		{itemSpellout, "komma"},
		{itemRightSub, " →→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

	//
	input = "←← komma →→;"
	exp = Result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
		{itemRightSub, " →→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

	//
	input = "←%cardinal-neuter← komma →%cardinal-reale→;"
	exp = Result{
		{itemLeftSub, "←%cardinal-neuter← "},
		{itemSpellout, "komma"},
		{itemRightSub, " →%cardinal-reale→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

}

func TestOptionalSub(t *testing.T) {
	var input string
	var exp Result
	var l *Lexer

	//
	input = "[←← ]komma →→;"
	exp = Result{
		{itemLeftBracket, "["},
		{itemLeftSub, "←← "},
		{itemRightBracket, "]"},
		{itemSpellout, "komma"},
		{itemRightSub, " →→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

	//
	input = "←← komma[ →→];"
	exp = Result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
		{itemLeftBracket, "["},
		{itemRightSub, " →→"},
		{itemRightBracket, "]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}

	//
	input = "←← komma[ →%cardinal-reale→];"
	exp = Result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
		{itemLeftBracket, "["},
		{itemRightSub, " →%cardinal-reale→"},
		{itemRightBracket, "]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result) {
		t.Error(err)
	}
}
