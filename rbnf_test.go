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

	r = NewIntRule(10, "", "", "tio", "", "", 10)
	if w, g := 10, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(100, "", "", "hundra", "", "", 10)
	if w, g := 100, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(200, "", "", "hundra", "", "", 10)
	if w, g := 100, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(2000, "", "", "tusen", "", "", 10)
	if w, g := 1000, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(2000, "", "", "tusen", "", ">", 10)
	if w, g := 1000, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = NewIntRule(2000, "", "", "tusen", "", ">>", 10)
	if w, g := 1000, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	// <rbnfrule value="1100" radix="100">←←­hundra[­→→];</rbnfrule>

	r = NewIntRule(1100, "", "", "hundra", "", ">>", 100)
	if w, g := 100, r.Base.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

}

func Test_Spellout1(t *testing.T) {
	defaultRules := RuleSet{
		Name: "default",
		Rules: []BaseRule{
			NewIntRule(0, "", "", "noll", "", "", 10),
			NewIntRule(1, "", "", "ett", "", "", 10),
			NewIntRule(2, "", "", "två", "", "", 10),
			NewIntRule(3, "", "", "tre", "", "", 10),
			NewIntRule(4, "", "", "fyra", "", "", 10),
			NewIntRule(5, "", "", "fem", "", "", 10),
			NewIntRule(6, "", "", "sex", "", "", 10),
			NewIntRule(7, "", "", "sju", "", "", 10),
			NewIntRule(8, "", "", "åtta", "", "", 10),
			NewIntRule(9, "", "", "nio", "", "", 10),
			NewIntRule(10, "", "", "tio", "", "", 10),
			NewIntRule(11, "", "", "elva", "", "", 10),
			NewIntRule(12, "", "", "tolv", "", "", 10),
			NewIntRule(13, "", "", "tretton", "", "", 10),
			NewIntRule(14, "", "", "fjorton", "", "", 10),
			NewIntRule(15, "", "", "femton", "", "", 10),
			NewIntRule(16, "", "", "sexton", "", "", 10),
			NewIntRule(17, "", "", "sjutton", "", "", 10),
			NewIntRule(18, "", "", "arton", "", "", 10),
			NewIntRule(19, "", "", "nitton", "", "", 10),
			NewIntRule(20, "", "", "tjugo", "-", "[>>]", 10),
			NewIntRule(30, "", "", "trettio", "-", "[>>]", 10),
			NewIntRule(40, "", "", "fyrtio", "-", "[>>]", 10),
			NewIntRule(50, "", "", "femtio", "-", "[>>]", 10),
			NewIntRule(60, "", "", "sextio", "-", "[>>]", 10),
			NewIntRule(70, "", "", "sjuttio", "-", "[>>]", 10),
			NewIntRule(80, "", "", "åttio", "-", "[>>]", 10),
			NewIntRule(90, "", "", "nittio", "-", "[>>]", 10),

			NewIntRule(100, "<<", " ", "hundra", " ", "[>>]", 10),

			NewIntRule(1000, "", " ", "ettusen", " ", "[>>]", 10),
			NewIntRule(2000, "%spellout-cardinal-reale", " ", "tusen", " ", "[>>]", 10),

			NewIntRule(1000000, "", " ", "en miljon", " ", "[>>]", 10),
			NewIntRule(2000000, "%spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10),
			NewIntRule(1000000000, "", "", "en miljard", " ", "[>>]", 10),
			NewIntRule(2000000000, "%spellout-cardinal-reale", " ", "miljarder", " ", "[>>]", 10),
		},
	}

	spelloutCardinalReale := RuleSet{
		Name: "spellout-cardinal-reale",
		Rules: []BaseRule{
			NewIntRule(0, "", "", "noll", "", "", 10),
			NewIntRule(1, "", "", "en", "", "", 10),
			NewIntRule(2, "", "", "två", "", "", 10),
			NewIntRule(3, "", "", "tre", "", "", 10),
			NewIntRule(4, "", "", "fyra", "", "", 10),
			NewIntRule(5, "", "", "fem", "", "", 10),
			NewIntRule(6, "", "", "sex", "", "", 10),
			NewIntRule(7, "", "", "sju", "", "", 10),
			NewIntRule(8, "", "", "åtta", "", "", 10),
			NewIntRule(9, "", "", "nio", "", "", 10),
			NewIntRule(10, "", "", "tio", "", "", 10),
			NewIntRule(11, "", "", "elva", "", "", 10),
			NewIntRule(12, "", "", "tolv", "", "", 10),
			NewIntRule(13, "", "", "tretton", "", "", 10),
			NewIntRule(14, "", "", "fjorton", "", "", 10),
			NewIntRule(15, "", "", "femton", "", "", 10),
			NewIntRule(16, "", "", "sexton", "", "", 10),
			NewIntRule(17, "", "", "sjutton", "", "", 10),
			NewIntRule(18, "", "", "arton", "", "", 10),
			NewIntRule(19, "", "", "nitton", "", "", 10),
			NewIntRule(20, "", "", "tjugo", "-", "[>>]", 10),
			NewIntRule(30, "", "", "trettio", "-", "[>>]", 10),
			NewIntRule(40, "", "", "fyrtio", "-", "[>>]", 10),
			NewIntRule(50, "", "", "femtio", "-", "[>>]", 10),
			NewIntRule(60, "", "", "sextio", "-", "[>>]", 10),
			NewIntRule(70, "", "", "sjuttio", "-", "[>>]", 10),
			NewIntRule(80, "", "", "åttio", "-", "[>>]", 10),
			NewIntRule(90, "", "", "nittio", "-", "[>>]", 10),
			NewIntRule(100, "%spellout-cardinal-neuter", " ", "hundra", " ", "[>>]", 10),
			NewIntRule(1000, "", " ", "ettusen", "-", "[>>]", 10),
			NewIntRule(2000, "%spellout-cardinal-reale", " ", "tusen", " ", "[>>]", 10),
			NewIntRule(1000000, "", " ", "en miljon", " ", "[>>]", 10),
			NewIntRule(2000000, "%spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10),
		},
	}

	spelloutCardinalNeuter := RuleSet{
		Name: "spellout-cardinal-neuter",
		Rules: []BaseRule{
			NewIntRule(0, "", "", "noll", "", "", 10),
			NewIntRule(1, "", "", "ett", "", "", 10),
			NewIntRule(2, "", "", "två", "", "", 10),
			NewIntRule(3, "", "", "tre", "", "", 10),
			NewIntRule(4, "", "", "fyra", "", "", 10),
			NewIntRule(5, "", "", "fem", "", "", 10),
			NewIntRule(6, "", "", "sex", "", "", 10),
			NewIntRule(7, "", "", "sju", "", "", 10),
			NewIntRule(8, "", "", "åtta", "", "", 10),
			NewIntRule(9, "", "", "nio", "", "", 10),
			NewIntRule(10, "", "", "tio", "", "", 10),
			NewIntRule(11, "", "", "elva", "", "", 10),
			NewIntRule(12, "", "", "tolv", "", "", 10),
			NewIntRule(13, "", "", "tretton", "", "", 10),
			NewIntRule(14, "", "", "fjorton", "", "", 10),
			NewIntRule(15, "", "", "femton", "", "", 10),
			NewIntRule(16, "", "", "sexton", "", "", 10),
			NewIntRule(17, "", "", "sjutton", "", "", 10),
			NewIntRule(18, "", "", "arton", "", "", 10),
			NewIntRule(19, "", "", "nitton", "", "", 10),
			NewIntRule(20, "", "", "tjugo", "-", "[>>]", 10),
			NewIntRule(30, "", "", "trettio", "-", "[>>]", 10),
			NewIntRule(40, "", "", "fyrtio", "-", "[>>]", 10),
			NewIntRule(50, "", "", "femtio", "-", "[>>]", 10),
			NewIntRule(60, "", "", "sextio", "-", "[>>]", 10),
			NewIntRule(70, "", "", "sjuttio", "-", "[>>]", 10),
			NewIntRule(80, "", "", "åttio", "-", "[>>]", 10),
			NewIntRule(90, "", "", "nittio", "-", "[>>]", 10),
			NewIntRule(100, "%spellout-cardinal-neuter", "", "hundra", " ", "[>>]", 10),
			NewIntRule(1000, "", "", "ettusen", " ", "[>>]", 10),
			NewIntRule(2000, "%spellout-cardinal-reale", "", "tusen", " ", "[>>]", 10),
			NewIntRule(1000000, "", "", "en miljon", " ", "[>>]", 10),
			NewIntRule(2000000, "%spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10),
			NewIntRule(1000000000, "", "", "en miljard", " ", "[>>]", 10),
			NewIntRule(2000000000, "%spellout-cardinal-reale", " ", "miljarder", " ", "[>>]", 10),
		},
	}

	g, err := NewRuleSetGroup(
		"spellout-cardinal",
		[]RuleSet{
			defaultRules,
			spelloutCardinalReale,
			spelloutCardinalNeuter,
		})
	if err != nil {
		t.Errorf("Couldn't create rule set group : %v", err)
	}

	// TEST
	var exp, res string

	res, err = g.Spellout("12", "default")
	exp = "tolv"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("3106", "default")
	exp = "tre tusen ett hundra sex"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("725601", "default")
	exp = "sju hundra tjugo-fem tusen sex hundra ett"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("681", "default")
	exp = "sex hundra åttio-ett"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2000000", "default")
	exp = "två miljoner"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2510000", "default")
	exp = "två miljoner fem hundra tio tusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2500000", "default")
	exp = "två miljoner fem hundra tusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("4123000", "default")
	exp = "fyra miljoner ett hundra tjugo-tre tusen"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("31607106", "default")
	exp = "trettio-en miljoner sex hundra sju tusen ett hundra sex"
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
			NewStringRule("-x", "", "", "minus", " ", ">>"),
			NewStringRule("x.x", "<<", " ", "komma", " ", ">>"),
			NewIntRule(0, "", "", "noll", "", "", 10),
			NewIntRule(1, "", "", "ett", "", "", 10),
			NewIntRule(2, "", "", "två", "", "", 10),
			NewIntRule(3, "", "", "tre", "", "", 10),
			NewIntRule(4, "", "", "fyra", "", "", 10),
			NewIntRule(5, "", "", "fem", "", "", 10),
			NewIntRule(6, "", "", "sex", "", "", 10),
			NewIntRule(7, "", "", "sju", "", "", 10),
			NewIntRule(8, "", "", "åtta", "", "", 10),
			NewIntRule(9, "", "", "nio", "", "", 10),
			NewIntRule(10, "", "", "tio", "", "", 10),
			NewIntRule(11, "", "", "elva", "", "", 10),
			NewIntRule(12, "", "", "tolv", "", "", 10),
			NewIntRule(13, "", "", "tretton", "", "", 10),
			NewIntRule(14, "", "", "fjorton", "", "", 10),
			NewIntRule(15, "", "", "femton", "", "", 10),
			NewIntRule(16, "", "", "sexton", "", "", 10),
			NewIntRule(17, "", "", "sjutton", "", "", 10),
			NewIntRule(18, "", "", "arton", "", "", 10),
			NewIntRule(19, "", "", "nitton", "", "", 10),
			NewIntRule(20, "", "", "tjugo", "-", "[>>]", 10),
			NewIntRule(30, "", "", "trettio", "-", "[>>]", 10),
			NewIntRule(40, "", "", "fyrtio", "-", "[>>]", 10),
			NewIntRule(50, "", "", "femtio", "-", "[>>]", 10),
			NewIntRule(60, "", "", "sextio", "-", "[>>]", 10),
			NewIntRule(70, "", "", "sjuttio", "-", "[>>]", 10),
			NewIntRule(80, "", "", "åttio", "-", "[>>]", 10),
			NewIntRule(90, "", "", "nittio", "-", "[>>]", 10),
			NewIntRule(100, "<<", "", "hundra", " ", "[>>]", 10),
			NewIntRule(1100, "<<", " ", "hundra", " ", "[>>]", 100),
			NewIntRule(2000, "<<", " ", "tusen", " ", "[>>]", 10),
		},
	}

	g, err := NewRuleSetGroup(
		"years",
		[]RuleSet{
			defaultRules,
		})
	if err != nil {
		t.Errorf("Couldn't create rule set group : %v", err)
	}

	// TEST
	var exp, res string

	res, err = g.Spellout("12", "default")
	exp = "tolv"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("1803", "default")
	exp = "arton hundra tre"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("1983", "default")
	exp = "nitton hundra åttio-tre"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("2001", "default")
	exp = "två tusen ett"
	if err != nil {
		t.Error(err)
	}
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Spellout("-2001x", "default")
	exp = "No matching base rule for"
	if err == nil {
		t.Errorf("Expected error, found %v", err)
	}
	if !strings.Contains(err.Error(), exp) {
		t.Errorf(fs, exp, err)
	}

	res, err = g.Spellout("-2001", "default")
	exp = "minus två tusen ett"
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
