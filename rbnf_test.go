package rbnf

import (
	"strings"
	"testing"
)

var fs = "Expected '%v', got '%v'"

func Test_exp(t *testing.T) {
	if w, g := 10, exp(10, 1); w != g {
		t.Errorf(fs, w, g)
	}
	if w, g := 100, exp(10, 2); w != g {
		t.Errorf(fs, w, g)
	}
	if w, g := 1, exp(10, 0); w != g {
		t.Errorf(fs, w, g)
	}
}

func Test_Divisor(t *testing.T) {
	var r BaseRule

	r = NewIntRule(10, 10, "tio")
	if w, g := 10, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(100, 10, "hundra")
	if w, g := 100, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(200, 10, "hundra")
	if w, g := 100, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(2000, 10, "tusen", "")
	if w, g := 1000, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(2000, 10, "tusen", ">>")
	if w, g := 1000, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(1100, 100, "hundra", ">>")
	if w, g := 100, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

}

func Test_Spellout1(t *testing.T) {
	defaultRules := RuleSet{
		Name: "default",
		Rules: []BaseRule{
			NewIntRule(0, 10, "noll"),
			NewIntRule(1, 10, "ett"),
			NewIntRule(2, 10, "två"),
			NewIntRule(3, 10, "tre"),
			NewIntRule(4, 10, "fyra"),
			NewIntRule(5, 10, "fem"),
			NewIntRule(6, 10, "sex"),
			NewIntRule(7, 10, "sju"),
			NewIntRule(8, 10, "åtta"),
			NewIntRule(9, 10, "nio"),
			NewIntRule(10, 10, "tio"),
			NewIntRule(11, 10, "elva"),
			NewIntRule(12, 10, "tolv"),
			NewIntRule(13, 10, "tretton"),
			NewIntRule(14, 10, "fjorton"),
			NewIntRule(15, 10, "femton"),
			NewIntRule(16, 10, "sexton"),
			NewIntRule(17, 10, "sjutton"),
			NewIntRule(18, 10, "arton"),
			NewIntRule(19, 10, "nitton"),
			NewIntRule(20, 10, "tjugo", "[-]", "[>>]"),
			NewIntRule(30, 10, "trettio", "[-]", "[>>]"),
			NewIntRule(40, 10, "fyrtio", "[-]", "[>>]"),
			NewIntRule(50, 10, "femtio", "[-]", "[>>]"),
			NewIntRule(60, 10, "sextio", "[-]", "[>>]"),
			NewIntRule(70, 10, "sjuttio", "[-]", "[>>]"),
			NewIntRule(80, 10, "åttio", "[-]", "[>>]"),
			NewIntRule(90, 10, "nittio", "[-]", "[>>]"),

			NewIntRule(100, 10, "<<", " ", "hundra", "[ ]", "[>>]"),

			NewIntRule(1000, 10, " ", "ettusen", "[ ]", "[>>]"),
			NewIntRule(2000, 10, "<%spellout-cardinal-reale<", " ", "tusen", "[ ]", "[>>]"),

			NewIntRule(1000000, 10, " ", "en miljon", "[ ]", "[>>]"),
			NewIntRule(2000000, 10, "<%spellout-cardinal-reale<", " ", "miljoner", "[ ]", "[>>]"),
			NewIntRule(1000000000, 10, "en miljard", "[ ]", "[>>]"),
			NewIntRule(2000000000, 10, "<%spellout-cardinal-reale<", " ", "miljarder", "[ ]", "[>>]"),
		},
	}

	spelloutCardinalReale := RuleSet{
		Name: "spellout-cardinal-reale",
		Rules: []BaseRule{
			NewIntRule(0, 10, "noll"),
			NewIntRule(1, 10, "en"),
			NewIntRule(2, 10, "två"),
			NewIntRule(3, 10, "tre"),
			NewIntRule(4, 10, "fyra"),
			NewIntRule(5, 10, "fem"),
			NewIntRule(6, 10, "sex"),
			NewIntRule(7, 10, "sju"),
			NewIntRule(8, 10, "åtta"),
			NewIntRule(9, 10, "nio"),
			NewIntRule(10, 10, "tio"),
			NewIntRule(11, 10, "elva"),
			NewIntRule(12, 10, "tolv"),
			NewIntRule(13, 10, "tretton"),
			NewIntRule(14, 10, "fjorton"),
			NewIntRule(15, 10, "femton"),
			NewIntRule(16, 10, "sexton"),
			NewIntRule(17, 10, "sjutton"),
			NewIntRule(18, 10, "arton"),
			NewIntRule(19, 10, "nitton"),
			NewIntRule(20, 10, "tjugo", "[-]", "[>>]"),
			NewIntRule(30, 10, "trettio", "[-]", "[>>]"),
			NewIntRule(40, 10, "fyrtio", "[-]", "[>>]"),
			NewIntRule(50, 10, "femtio", "[-]", "[>>]"),
			NewIntRule(60, 10, "sextio", "[-]", "[>>]"),
			NewIntRule(70, 10, "sjuttio", "[-]", "[>>]"),
			NewIntRule(80, 10, "åttio", "[-]", "[>>]"),
			NewIntRule(90, 10, "nittio", "[-]", "[>>]"),
			NewIntRule(100, 10, "<%spellout-cardinal-neuter<", " ", "hundra", "[ ]", "[>>]"),
			NewIntRule(1000, 10, " ", "ettusen", "[-]", "[>>]"),
			NewIntRule(2000, 10, "<%spellout-cardinal-reale<", " ", "tusen", "[ ]", "[>>]"),
			NewIntRule(1000000, 10, " ", "en miljon", "[ ]", "[>>]"),
			NewIntRule(2000000, 10, "<%spellout-cardinal-reale<", " ", "miljoner", "[ ]", "[>>]"),
		},
	}

	spelloutCardinalNeuter := RuleSet{
		Name: "spellout-cardinal-neuter",
		Rules: []BaseRule{
			NewIntRule(0, 10, "noll"),
			NewIntRule(1, 10, "ett"),
			NewIntRule(2, 10, "två"),
			NewIntRule(3, 10, "tre"),
			NewIntRule(4, 10, "fyra"),
			NewIntRule(5, 10, "fem"),
			NewIntRule(6, 10, "sex"),
			NewIntRule(7, 10, "sju"),
			NewIntRule(8, 10, "åtta"),
			NewIntRule(9, 10, "nio"),
			NewIntRule(10, 10, "tio"),
			NewIntRule(11, 10, "elva"),
			NewIntRule(12, 10, "tolv"),
			NewIntRule(13, 10, "tretton"),
			NewIntRule(14, 10, "fjorton"),
			NewIntRule(15, 10, "femton"),
			NewIntRule(16, 10, "sexton"),
			NewIntRule(17, 10, "sjutton"),
			NewIntRule(18, 10, "arton"),
			NewIntRule(19, 10, "nitton"),
			NewIntRule(20, 10, "tjugo", "[-]", "[>>]"),
			NewIntRule(30, 10, "trettio", "[-]", "[>>]"),
			NewIntRule(40, 10, "fyrtio", "[-]", "[>>]"),
			NewIntRule(50, 10, "femtio", "[-]", "[>>]"),
			NewIntRule(60, 10, "sextio", "[-]", "[>>]"),
			NewIntRule(70, 10, "sjuttio", "[-]", "[>>]"),
			NewIntRule(80, 10, "åttio", "[-]", "[>>]"),
			NewIntRule(90, 10, "nittio", "[-]", "[>>]"),
			NewIntRule(100, 10, "<%spellout-cardinal-neuter<", "hundra", "[ ]", "[>>]"),
			NewIntRule(1000, 10, "ettusen", "[ ]", "[>>]"),
			NewIntRule(2000, 10, "<%spellout-cardinal-reale<", "tusen", "[ ]", "[>>]"),
			NewIntRule(1000000, 10, "en miljon", "[ ]", "[>>]"),
			NewIntRule(2000000, 10, "<%spellout-cardinal-reale<", " ", "miljoner", "[ ]", "[>>]"),
			NewIntRule(1000000000, 10, "en miljard", "[ ]", "[>>]"),
			NewIntRule(2000000000, 10, "<%spellout-cardinal-reale<", " ", "miljarder", "[ ]", "[>>]"),
		},
	}
	spelloutCardinalNeuter2 := RuleSet{
		Name: "spellout-cardinal-neuter-2",
		Rules: []BaseRule{
			NewIntRule(0, 10, "=%spellout-cardinal-neuter="),
		},
	}
	g, err := NewRuleSetGroup(
		"spellout-cardinal",
		"sv",
		[]RuleSet{
			defaultRules,
			spelloutCardinalReale,
			spelloutCardinalNeuter,
			spelloutCardinalNeuter2,
		})
	if err != nil {
		t.Errorf("Couldn't create rule set group : %v", err)
	}

	// TEST
	var exp, res string

	res, err = g.Spellout("12", "default", false)
	exp = "tolv"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("3106", "default", false)
	exp = "tre tusen ett hundra sex"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("725601", "default", false)
	exp = "sju hundra tjugo-fem tusen sex hundra ett"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("681", "default", false)
	exp = "sex hundra åttio-ett"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("20000", "default", false)
	exp = "tjugo tusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2000000", "default", false)
	exp = "två miljoner"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("20", "default", false)
	exp = "tjugo"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("20000000", "default", false)
	exp = "tjugo miljoner"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("200000000", "default", false)
	exp = "två hundra miljoner"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2510000", "default", false)
	exp = "två miljoner fem hundra tio tusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2500000", "default", false)
	exp = "två miljoner fem hundra tusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2001000", "default", false)
	exp = "två miljoner ettusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("4123000", "default", false)
	exp = "fyra miljoner ett hundra tjugo-tre tusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("31607106", "default", false)
	exp = "trettio-en miljoner sex hundra sju tusen ett hundra sex"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("0", "default", false)
	exp = "noll"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}
}

func Test_Spellout2(t *testing.T) {
	defaultRules := RuleSet{
		Name: "default",
		Rules: []BaseRule{
			NewStringRule("-x", "minus", " ", ">>"),
			NewStringRule("x.x", "<<", " ", "komma", " ", ">>"),
			NewIntRule(0, 10, "noll"),
			NewIntRule(1, 10, "ett"),
			NewIntRule(2, 10, "två"),
			NewIntRule(3, 10, "tre"),
			NewIntRule(4, 10, "fyra"),
			NewIntRule(5, 10, "fem"),
			NewIntRule(6, 10, "sex"),
			NewIntRule(7, 10, "sju"),
			NewIntRule(8, 10, "åtta"),
			NewIntRule(9, 10, "nio"),
			NewIntRule(10, 10, "tio"),
			NewIntRule(11, 10, "elva"),
			NewIntRule(12, 10, "tolv"),
			NewIntRule(13, 10, "tretton"),
			NewIntRule(14, 10, "fjorton"),
			NewIntRule(15, 10, "femton"),
			NewIntRule(16, 10, "sexton"),
			NewIntRule(17, 10, "sjutton"),
			NewIntRule(18, 10, "arton"),
			NewIntRule(19, 10, "nitton"),
			NewIntRule(20, 10, "tjugo", "[-]", "[>>]"),
			NewIntRule(30, 10, "trettio", "[-]", "[>>]"),
			NewIntRule(40, 10, "fyrtio", "[-]", "[>>]"),
			NewIntRule(50, 10, "femtio", "[-]", "[>>]"),
			NewIntRule(60, 10, "sextio", "[-]", "[>>]"),
			NewIntRule(70, 10, "sjuttio", "[-]", "[>>]"),
			NewIntRule(80, 10, "åttio", "[-]", "[>>]"),
			NewIntRule(90, 10, "nittio", "[-]", "[>>]"),
			NewIntRule(100, 10, "<<", "hundra", "[ ]", "[>>]"),
			NewIntRule(1100, 100, "<<", " ", "hundra", "[ ]", "[>>]"),
			NewIntRule(2000, 10, "<<", " ", "tusen", "[ ]", "[>>]"),
		},
	}

	g, err := NewRuleSetGroup(
		"years",
		"sv",
		[]RuleSet{
			defaultRules,
		})
	if err != nil {
		t.Errorf("Couldn't create rule set group : %v", err)
	}

	// TEST
	var exp, res string

	res, err = g.Spellout("12", "default", false)
	exp = "tolv"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("1803", "default", false)
	exp = "arton hundra tre"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("1983", "default", false)
	exp = "nitton hundra åttio-tre"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2001", "default", false)
	exp = "två tusen ett"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("-2001x", "default", false)
	exp = "No matching base rule for"
	if err == nil {
		t.Errorf("Expected error, found %v", err)
	}
	if !strings.Contains(err.Error(), exp) {
		t.Errorf(fs, exp, err)
	}

	res, err = g.Spellout("-2001", "default", false)
	exp = "minus två tusen ett"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

}

func Test_SpelloutOperationEqSign(t *testing.T) {
	defaultRules := RuleSet{
		Name: "default",
		Rules: []BaseRule{
			NewStringRule("-x", "minus", " ", ">>"),
			NewStringRule("x.x", "<<", " ", "komma", " ", ">>"),
			NewIntRule(0, 10, "noll"),
			NewIntRule(1, 10, "ett"),
			NewIntRule(2, 10, "två"),
			NewIntRule(3, 10, "tre"),
			NewIntRule(4, 10, "fyra"),
			NewIntRule(5, 10, "fem"),
			NewIntRule(6, 10, "sex"),
			NewIntRule(7, 10, "sju"),
			NewIntRule(8, 10, "åtta"),
			NewIntRule(9, 10, "nio"),
			NewIntRule(10, 10, "tio"),
			NewIntRule(11, 10, "elva"),
			NewIntRule(12, 10, "tolv"),
			NewIntRule(13, 10, "tretton"),
			NewIntRule(14, 10, "fjorton"),
			NewIntRule(15, 10, "femton"),
			NewIntRule(16, 10, "sexton"),
			NewIntRule(17, 10, "sjutton"),
			NewIntRule(18, 10, "arton"),
			NewIntRule(19, 10, "nitton"),
			NewIntRule(20, 10, "tjugo", "[-]", "[>>]"),
			NewIntRule(30, 10, "trettio", "[-]", "[>>]"),
			NewIntRule(40, 10, "fyrtio", "[-]", "[>>]"),
			NewIntRule(50, 10, "femtio", "[-]", "[>>]"),
			NewIntRule(60, 10, "sextio", "[-]", "[>>]"),
			NewIntRule(70, 10, "sjuttio", "[-]", "[>>]"),
			NewIntRule(80, 10, "åttio", "[-]", "[>>]"),
			NewIntRule(90, 10, "nittio", "[-]", "[>>]"),
			NewIntRule(100, 10, "<<", "hundra", "[ ]", "[>>]"),
			NewIntRule(1100, 100, "<<", " ", "hundra", "[ ]", "[>>]"),
			NewIntRule(2000, 10, "<<", " ", "tusen", "[ ]", "[>>]"),
		},
	}
	rules2 := RuleSet{
		Name: "rules2",
		Rules: []BaseRule{
			NewIntRule(0, 10, "=%default="),
		},
	}

	g, err := NewRuleSetGroup(
		"years",
		"sv",
		[]RuleSet{
			defaultRules,
			rules2,
		})
	if err != nil {
		t.Errorf("Couldn't create rule set group : %v", err)
	}

	// TEST
	var exp, res string

	res, err = g.Spellout("12", "rules2", false)
	exp = "tolv"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

}

func Test_StringMatch(t *testing.T) {
	var r BaseRule
	var exp, res MatchResult
	var ok bool
	var in string

	r = NewStringRule("-x", "", "", "minus", " ", ">>")
	in = "-18"
	res, ok = r.Match(in)
	exp = MatchResult{ForwardLeft: "", ForwardRight: "18"}
	if !ok {
		t.Errorf("No match result for %s", in)
	}
	if exp != res {
		t.Errorf(fs, exp, res)
	}

	r = NewStringRule("x.x", "<<", " ", "komma", " ", ">>")
	in = "3.18"
	res, ok = r.Match(in)
	exp = MatchResult{ForwardLeft: "3", ForwardRight: "18"}
	if !ok {
		t.Errorf("No match result for %s", in)
	}
	if exp != res {
		t.Errorf(fs, exp, res)
	}

	r = NewStringRule("x,x", "<<", " ", "komma", " ", ">>")
	in = "3,18"
	res, ok = r.Match(in)
	exp = MatchResult{ForwardLeft: "3", ForwardRight: "18"}
	if !ok {
		t.Errorf("No match result for %s", in)
	}
	if exp != res {
		t.Errorf(fs, exp, res)
	}

	r = NewStringRule("x%", "<<", " ", "procent", " ", "")
	in = "316%"
	res, ok = r.Match(in)
	exp = MatchResult{ForwardLeft: "316", ForwardRight: ""}
	if !ok {
		t.Errorf("No match result for %s", in)
	}
	if exp != res {
		t.Errorf(fs, exp, res)
	}
}

func Test_SpelloutDE(t *testing.T) {
	defaultRules := RuleSet{
		Name: "spellout-numbering",
		Rules: []BaseRule{
			NewStringRule("-x", "minus", " ", ">>"),
			NewStringRule("x.x", "<<", " ", "komma", " ", ">>"),
			NewIntRule(0, 10, "null"),
			NewIntRule(1, 10, "eins"),
			NewIntRule(2, 10, "zwei"),
			NewIntRule(3, 10, "drei"),
			NewIntRule(4, 10, "vier"),
			NewIntRule(5, 10, "fünf"),
			NewIntRule(6, 10, "sechs"),
			NewIntRule(7, 10, "sieben"),
			NewIntRule(8, 10, "acht"),
			NewIntRule(9, 10, "neun"),
			NewIntRule(10, 10, "zehn"),
			NewIntRule(11, 10, "elf"),
			NewIntRule(12, 10, "zwölf"),
			NewIntRule(13, 10, ">>zehn"),
			NewIntRule(16, 10, "sechzehn"),
			NewIntRule(17, 10, "siebzehn"),
			NewIntRule(18, 10, ">>zehn"),
			NewIntRule(20, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "zwanzig"),
			NewIntRule(30, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "dreißig"),
			NewIntRule(40, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "vierzig"),
			NewIntRule(50, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "fünfzig"),
			NewIntRule(60, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "sechzig"),
			NewIntRule(70, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "siebzig"),
			NewIntRule(80, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "achtzig"),
			NewIntRule(90, 10, "[>%spellout-cardinal-masculine>]", "[-und-]", "neunzig"),
			NewIntRule(100, 10, "ERROR"),
			NewIntRule(1000000000000000, 10, "=#,##0="),
			NewIntRule(1000000000000000000, 10, "=0="),
		},
	}
	spelloutCardinalMasculine := RuleSet{
		Name: "spellout-cardinal-masculine",
		Rules: []BaseRule{
			NewStringRule("-x", "minus", " ", ">>"),
			NewStringRule("x.x", "<<", " ", "komma", " ", ">>"),
			NewIntRule(0, 10, "null"),
			NewIntRule(1, 10, "ein"),
			NewIntRule(2, 10, "=%spellout-numbering="),
		},
	}

	g, err := NewRuleSetGroup(
		"default",
		"de",
		[]RuleSet{
			defaultRules,
			spelloutCardinalMasculine,
		})
	if err != nil {
		t.Errorf("Couldn't create rule set group : %v", err)
	}

	// TEST
	var exp, res string

	//
	res, err = g.Spellout("12", "spellout-numbering", false)
	exp = "zwölf"
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	res, err = g.Spellout("45", "spellout-numbering", false)
	exp = "fünf-und-vierzig"
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	res, err = g.Spellout("100", "spellout-numbering", false)
	if err == nil {
		t.Error("expected error here")
	}

	//
	res, err = g.Spellout("1000000000000000", "spellout-numbering", false)
	exp = "1.000.000.000.000.000"
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	res, err = g.Spellout("1000000000000000000", "spellout-numbering", false)
	exp = "1.000.000.000.000.000.000"
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

}

func TestHashFormat(t *testing.T) {

	var lang, input, res, exp string
	var err error

	//
	lang = "en"
	input = "1000000000000000000"
	exp = "1,000,000,000,000,000,000"
	res, err = format(input, lang, "#,##", false)
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	lang = "en"
	input = "12000.3789"
	exp = "12,000.3789"
	res, err = format(input, lang, "#,##", false)
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	lang = "de"
	input = "1000000000000000000"
	exp = "1.000.000.000.000.000.000"
	res, err = format(input, lang, "#,##", false)
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	lang = "de"
	input = "12000.3789"
	exp = "12.000,3789"
	res, err = format(input, lang, "#,##", false)
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	lang = "sv"
	input = "1000000000000000000"
	exp = "1 000 000 000 000 000 000" // non-breaking space \u00A0
	res, err = format(input, lang, "#,##", false)
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	lang = "sv"
	input = "12000.3789"
	exp = "12 000,3789" // non-breaking space \u00A0
	res, err = format(input, lang, "#,##", false)
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}

	//
	lang = "bn"
	input = "123456.78"
	exp = "১,২৩,৪৫৬.৭৮"
	res, err = format(input, lang, "#,##", false)
	if err != nil {
		t.Error(err)
	} else if res != exp {
		t.Errorf(fs, exp, res)
	}
}
