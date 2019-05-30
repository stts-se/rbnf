package main

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

func compare(input string, exp, got result) []error {
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
	var exp result
	var l *lexer
	var prematureEOIItem = item{itemError, "premature end of input"}

	//
	input = ""
	exp = result{prematureEOIItem}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "minus;"
	exp = result{{itemSpellout, "minus"}}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "minus"
	exp = result{{itemSpellout, "minus"}, prematureEOIItem}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

}

func TestSub(t *testing.T) {
	var input string
	var exp result
	var l *lexer

	//
	input = "←← komma;"
	exp = result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
	}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "komma →→;"
	exp = result{
		{itemSpellout, "komma"},
		{itemRightSub, " →→"},
	}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "←← komma →→;"
	exp = result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
		{itemRightSub, " →→"},
	}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "←%cardinal-neuter← komma →%cardinal-reale→;"
	exp = result{
		{itemLeftSub, "←%cardinal-neuter← "},
		{itemSpellout, "komma"},
		{itemRightSub, " →%cardinal-reale→"},
	}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

}

func TestOptionalSub(t *testing.T) {
	var input string
	var exp result
	var l *lexer

	//
	input = "[←← ]komma →→;"
	exp = result{
		{itemLeftBracket, "["},
		{itemLeftSub, "←← "},
		{itemRightBracket, "]"},
		{itemSpellout, "komma"},
		{itemRightSub, " →→"},
	}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "←← komma[ →→];"
	exp = result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
		{itemLeftBracket, "["},
		{itemRightSub, " →→"},
		{itemRightBracket, "]"},
	}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "←← komma[ →%cardinal-reale→];"
	exp = result{
		{itemLeftSub, "←← "},
		{itemSpellout, "komma"},
		{itemLeftBracket, "["},
		{itemRightSub, " →%cardinal-reale→"},
		{itemRightBracket, "]"},
	}
	l = lex(input)
	l.run()
	for _, err := range compare(input, exp, l.result) {
		t.Error(err)
	}
}
