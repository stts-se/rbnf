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

var fs = "For input '%s', expected '%#v', got '%#v'"
var fsindex = "For input '%s' item %d, expected '%s', got '%s'"

func compareResult(input string, exp, got result) []error {
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

func compareStrings(input string, exp, got []string) []error {
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
	var l *Lexer
	var prematureEOIItem = item{itemError, "premature end of input"}

	//
	input = ""
	exp = result{prematureEOIItem}
	l = Lex(input)
	l.Run()
	for _, err := range compareResult(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "minus"
	exp = result{{itemSub, "minus"}, prematureEOIItem}
	l = Lex(input)
	l.Run()
	for _, err := range compareResult(input, exp, l.result) {
		t.Error(err)
	}

	//
	input = "minus;"
	exp = result{{itemSub, "minus"}}
	l = Lex(input)
	l.Run()
	for _, err := range compareResult(input, exp, l.result) {
		t.Error(err)
	}

}

func TestSub(t *testing.T) {
	var input string
	var exp []string
	var l *Lexer

	//
	input = "←← komma;"
	exp = []string{
		"←←",
		" komma",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "komma →→;"
	exp = []string{
		"komma ",
		"→→",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←← komma →→;"
	exp = []string{
		"←←",
		" komma ",
		"→→",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←%cardinal-neuter← komma →%cardinal-reale→;"
	exp = []string{
		"←%cardinal-neuter←",
		" komma ",
		"→%cardinal-reale→",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

}

func TestOptionalSub(t *testing.T) {
	var input string
	var exp []string
	var l *Lexer

	//
	input = "[←← ]komma →→;"
	exp = []string{
		"[←←]",
		"[ ]",
		"komma ",
		"→→",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←← komma[ →→];"
	exp = []string{
		"←←",
		" komma",
		"[ ]",
		"[→→]",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "komma[ →→];"
	exp = []string{
		"komma",
		"[ ]",
		"[→→]",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←← komma[ →%cardinal-reale→];"
	exp = []string{
		"←←",
		" komma",
		"[ ]",
		"[→%cardinal-reale→]",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "en miljon[→→];"
	exp = []string{
		"en miljon",
		"[→→]",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "=%spellout-cardinal-neuter=de;"
	exp = []string{
		"=%spellout-cardinal-neuter=",
		"de",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "er =%spellout-cardinal-neuter= de;"
	exp = []string{
		"er ",
		"=%spellout-cardinal-neuter=",
		" de",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "=%spellout-numbering= miljoner tusen;"
	exp = []string{
		"=%spellout-numbering=",
		" miljoner tusen",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "­=%spellout-ordinal-feminine=;"
	exp = []string{
		"­",
		"=%spellout-ordinal-feminine=",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "tjugo[­→→];"
	exp = []string{
		"tjugo",
		"[­]",
		"[→→]",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

}

func TestMiscLangs(t *testing.T) {
	var input string
	var exp []string
	var l *Lexer

	//
	input = "sesenta[ y →→];"
	exp = []string{
		"sesenta",
		"[ y ]",
		"[→→]",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "[→%spellout-cardinal-masculine→­und­]fünfzig;"
	exp = []string{
		"[→%spellout-cardinal-masculine→]",
		"[­und­]",
		"fünfzig",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "' =%spellout-cardinal-masculine=;"
	exp = []string{
		"' ",
		"=%spellout-cardinal-masculine=",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "பன்னிரண்டு;"
	exp = []string{
		"பன்னிரண்டு",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "அறுபது;"
	exp = []string{
		"அறுபது",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "எழுநூறு;"
	exp = []string{
		"எழுநூறு",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = ", =%spellout-cardinal-verbose=;"
	exp = []string{
		", ",
		"=%spellout-cardinal-verbose=",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

}

func TestError(t *testing.T) {
	var input string
	var exp []string
	var l *Lexer

	//
	input = "ERROR;"
	exp = []string{
		"ERROR",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

}

func TestHashes(t *testing.T) {
	var input string
	var exp []string
	var l *Lexer

	//
	input = "=#,##0.#=;"
	exp = []string{
		"=#,##0.#=",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = ">#,##0.#>;"
	exp = []string{
		">#,##0.#>",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "=0=;"
	exp = []string{
		"=0=",
	}
	l = Lex(input)
	l.Run()
	for _, err := range compareStrings(input, exp, l.Result()) {
		t.Error(err)
	}

}
