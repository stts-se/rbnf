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
	var prematureEOIItem = Item{ItemError, "premature end of input"}

	//
	input = ""
	exp = Result{prematureEOIItem}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "minus"
	exp = Result{{ItemSub, "minus"}, prematureEOIItem}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "minus;"
	exp = Result{{ItemSub, "minus"}}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
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
		{ItemSub, "←←"},
		{ItemSub, " komma"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "komma →→;"
	exp = Result{
		{ItemSub, "komma "},
		{ItemSub, "→→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←← komma →→;"
	exp = Result{
		{ItemSub, "←←"},
		{ItemSub, " komma "},
		{ItemSub, "→→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←%cardinal-neuter← komma →%cardinal-reale→;"
	exp = Result{
		{ItemSub, "←%cardinal-neuter←"},
		{ItemSub, " komma "},
		{ItemSub, "→%cardinal-reale→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
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
		{ItemSub, "[←←]"},
		{ItemSub, "[ ]"},
		{ItemSub, "komma "},
		{ItemSub, "→→"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←← komma[ →→];"
	exp = Result{
		{ItemSub, "←←"},
		{ItemSub, " komma"},
		{ItemSub, "[ ]"},
		{ItemSub, "[→→]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "komma[ →→];"
	exp = Result{
		{ItemSub, "komma"},
		{ItemSub, "[ ]"},
		{ItemSub, "[→→]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "←← komma[ →%cardinal-reale→];"
	exp = Result{
		{ItemSub, "←←"},
		{ItemSub, " komma"},
		{ItemSub, "[ ]"},
		{ItemSub, "[→%cardinal-reale→]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "en miljon[→→];"
	exp = Result{
		{ItemSub, "en miljon"},
		{ItemSub, "[→→]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "=%spellout-cardinal-neuter=de;"
	exp = Result{
		{ItemSub, "=%spellout-cardinal-neuter="},
		{ItemSub, "de"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "er =%spellout-cardinal-neuter= de;"
	exp = Result{
		{ItemSub, "er "},
		{ItemSub, "=%spellout-cardinal-neuter="},
		{ItemSub, " de"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "=%spellout-numbering= miljoner tusen;"
	exp = Result{
		{ItemSub, "=%spellout-numbering="},
		{ItemSub, " miljoner tusen"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "­=%spellout-ordinal-feminine=;"
	exp = Result{
		{ItemSub, "­"},
		{ItemSub, "=%spellout-ordinal-feminine="},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "tjugo[­→→];"
	exp = Result{
		{ItemSub, "tjugo"},
		{ItemSub, "[­]"},
		{ItemSub, "[→→]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

}

func TestMiscLangs(t *testing.T) {
	var input string
	var exp Result
	var l *Lexer

	//
	input = "sesenta[ y →→];"
	exp = Result{
		{ItemSub, "sesenta"},
		{ItemSub, "[ y ]"},
		{ItemSub, "[→→]"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "[→%spellout-cardinal-masculine→­und­]fünfzig;"
	exp = Result{
		{ItemSub, "[→%spellout-cardinal-masculine→]"},
		{ItemSub, "[­und­]"},
		{ItemSub, "fünfzig"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "' =%spellout-cardinal-masculine=;"
	exp = Result{
		{ItemSub, "' "},
		{ItemSub, "=%spellout-cardinal-masculine="},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "பன்னிரண்டு;"
	exp = Result{
		{ItemSub, "பன்னிரண்டு"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "அறுபது;"
	exp = Result{
		{ItemSub, "அறுபது"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = "எழுநூறு;"
	exp = Result{
		{ItemSub, "எழுநூறு"},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

	//
	input = ", =%spellout-cardinal-verbose=;"
	exp = Result{
		{ItemSub, ", "},
		{ItemSub, "=%spellout-cardinal-verbose="},
	}
	l = Lex(input)
	l.Run()
	for _, err := range compare(input, exp, l.Result()) {
		t.Error(err)
	}

}
