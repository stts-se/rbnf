package rbnf

import (
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

	r = BaseRule{10, "", "", "tio", "", "", 10}
	if w, g := 10, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{100, "", "", "hundra", "", "", 10}
	if w, g := 100, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{200, "", "", "hundra", "", "", 10}
	if w, g := 100, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{2000, "", "", "tusen", "", "", 10}
	if w, g := 1000, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{2000, "", "", "tusen", "", ">", 10}
	if w, g := 1000, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{2000, "", "", "tusen", "", ">>", 10}
	if w, g := 1000, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	// <rbnfrule value="1100" radix="100">←←­hundra[­→→];</rbnfrule>

	r = BaseRule{1100, "", "", "hundra", "", ">>", 100}
	if w, g := 100, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

}

func Test_Expand1(t *testing.T) {
	defaultRules := RuleSet{
		Name: "default",
		Rules: []BaseRule{
			{0, "", "", "noll", "", "", 10},
			{1, "", "", "ett", "", "", 10},
			{2, "", "", "två", "", "", 10},
			{3, "", "", "tre", "", "", 10},
			{4, "", "", "fyra", "", "", 10},
			{5, "", "", "fem", "", "", 10},
			{6, "", "", "sex", "", "", 10},
			{7, "", "", "sju", "", "", 10},
			{8, "", "", "åtta", "", "", 10},
			{9, "", "", "nio", "", "", 10},
			{10, "", "", "tio", "", "", 10},
			{11, "", "", "elva", "", "", 10},
			{12, "", "", "tolv", "", "", 10},
			{13, "", "", "tretton", "", "", 10},
			{14, "", "", "fjorton", "", "", 10},
			{15, "", "", "femton", "", "", 10},
			{16, "", "", "sexton", "", "", 10},
			{17, "", "", "sjutton", "", "", 10},
			{18, "", "", "arton", "", "", 10},
			{19, "", "", "nitton", "", "", 10},
			{20, "", "", "tjugo", "-", "[>>]", 10},
			{30, "", "", "trettio", "-", "[>>]", 10},
			{40, "", "", "fyrtio", "-", "[>>]", 10},
			{50, "", "", "femtio", "-", "[>>]", 10},
			{60, "", "", "sextio", "-", "[>>]", 10},
			{70, "", "", "sjuttio", "-", "[>>]", 10},
			{80, "", "", "åttio", "-", "[>>]", 10},
			{90, "", "", "nittio", "-", "[>>]", 10},

			{100, "<<", " ", "hundra", " ", "[>>]", 10},

			{1000, "", " ", "ettusen", " ", "[>>]", 10},
			{2000, "spellout-cardinal-reale", " ", "tusen", " ", "[>>]", 10},

			{1000000, "", " ", "en miljon", " ", "[>>]", 10},
			{2000000, "spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10},
			{1000000000, "", "", "en miljard", " ", "[>>]", 10},
			{2000000000, "spellout-cardinal-reale", " ", "miljarder", " ", "[>>]", 10},
		},
	}

	spelloutCardinalReale := RuleSet{
		Name: "spellout-cardinal-reale",
		Rules: []BaseRule{
			{0, "", "", "noll", "", "", 10},
			{1, "", "", "en", "", "", 10},
			{2, "", "", "två", "", "", 10},
			{3, "", "", "tre", "", "", 10},
			{4, "", "", "fyra", "", "", 10},
			{5, "", "", "fem", "", "", 10},
			{6, "", "", "sex", "", "", 10},
			{7, "", "", "sju", "", "", 10},
			{8, "", "", "åtta", "", "", 10},
			{9, "", "", "nio", "", "", 10},
			{10, "", "", "tio", "", "", 10},
			{11, "", "", "elva", "", "", 10},
			{12, "", "", "tolv", "", "", 10},
			{13, "", "", "tretton", "", "", 10},
			{14, "", "", "fjorton", "", "", 10},
			{15, "", "", "femton", "", "", 10},
			{16, "", "", "sexton", "", "", 10},
			{17, "", "", "sjutton", "", "", 10},
			{18, "", "", "arton", "", "", 10},
			{19, "", "", "nitton", "", "", 10},
			{20, "", "", "tjugo", "-", "[>>]", 10},
			{30, "", "", "trettio", "-", "[>>]", 10},
			{40, "", "", "fyrtio", "-", "[>>]", 10},
			{50, "", "", "femtio", "-", "[>>]", 10},
			{60, "", "", "sextio", "-", "[>>]", 10},
			{70, "", "", "sjuttio", "-", "[>>]", 10},
			{80, "", "", "åttio", "-", "[>>]", 10},
			{90, "", "", "nittio", "-", "[>>]", 10},
			{100, "spellout-cardinal-neuter", " ", "hundra", " ", "[>>]", 10},
			{1000, "", " ", "ettusen", "-", "[>>]", 10},
			{2000, "spellout-cardinal-reale", " ", "tusen", " ", "[>>]", 10},
			{1000000, "", " ", "en miljon", " ", "[>>]", 10},
			{2000000, "spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10},
		},
	}

	spelloutCardinalNeuter := RuleSet{
		Name: "spellout-cardinal-neuter",
		Rules: []BaseRule{
			{0, "", "", "noll", "", "", 10},
			{1, "", "", "ett", "", "", 10},
			{2, "", "", "två", "", "", 10},
			{3, "", "", "tre", "", "", 10},
			{4, "", "", "fyra", "", "", 10},
			{5, "", "", "fem", "", "", 10},
			{6, "", "", "sex", "", "", 10},
			{7, "", "", "sju", "", "", 10},
			{8, "", "", "åtta", "", "", 10},
			{9, "", "", "nio", "", "", 10},
			{10, "", "", "tio", "", "", 10},
			{11, "", "", "elva", "", "", 10},
			{12, "", "", "tolv", "", "", 10},
			{13, "", "", "tretton", "", "", 10},
			{14, "", "", "fjorton", "", "", 10},
			{15, "", "", "femton", "", "", 10},
			{16, "", "", "sexton", "", "", 10},
			{17, "", "", "sjutton", "", "", 10},
			{18, "", "", "arton", "", "", 10},
			{19, "", "", "nitton", "", "", 10},
			{20, "", "", "tjugo", "-", "[>>]", 10},
			{30, "", "", "trettio", "-", "[>>]", 10},
			{40, "", "", "fyrtio", "-", "[>>]", 10},
			{50, "", "", "femtio", "-", "[>>]", 10},
			{60, "", "", "sextio", "-", "[>>]", 10},
			{70, "", "", "sjuttio", "-", "[>>]", 10},
			{80, "", "", "åttio", "-", "[>>]", 10},
			{90, "", "", "nittio", "-", "[>>]", 10},
			{100, "spellout-cardinal-neuter", "", "hundra", " ", "[>>]", 10},
			{1000, "", "", "ettusen", " ", "[>>]", 10},
			{2000, "spellout-cardinal-reale", "", "tusen", " ", "[>>]", 10},
			{1000000, "", "", "en miljon", " ", "[>>]", 10},
			{2000000, "spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10},
			{1000000000, "", "", "en miljard", " ", "[>>]", 10},
			{2000000000, "spellout-cardinal-reale", " ", "miljarder", " ", "[>>]", 10},
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

	res, err = g.Expand(12, "default")
	exp = "tolv"
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Expand(3106, "default")
	exp = "tre tusen ett hundra sex"
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Expand(31607106, "default")
	exp = "trettio-en miljoner sex hundra sju tusen ett hundra sex"
	if res != exp {
		t.Errorf(fs, exp, res)
	}

}

func Test_Expand2(t *testing.T) {
	defaultRules := RuleSet{
		Name: "default",
		Rules: []BaseRule{
			{0, "", "", "noll", "", "", 10},
			{1, "", "", "ett", "", "", 10},
			{2, "", "", "två", "", "", 10},
			{3, "", "", "tre", "", "", 10},
			{4, "", "", "fyra", "", "", 10},
			{5, "", "", "fem", "", "", 10},
			{6, "", "", "sex", "", "", 10},
			{7, "", "", "sju", "", "", 10},
			{8, "", "", "åtta", "", "", 10},
			{9, "", "", "nio", "", "", 10},
			{10, "", "", "tio", "", "", 10},
			{11, "", "", "elva", "", "", 10},
			{12, "", "", "tolv", "", "", 10},
			{13, "", "", "tretton", "", "", 10},
			{14, "", "", "fjorton", "", "", 10},
			{15, "", "", "femton", "", "", 10},
			{16, "", "", "sexton", "", "", 10},
			{17, "", "", "sjutton", "", "", 10},
			{18, "", "", "arton", "", "", 10},
			{19, "", "", "nitton", "", "", 10},
			{20, "", "", "tjugo", "-", "[>>]", 10},
			{30, "", "", "trettio", "-", "[>>]", 10},
			{40, "", "", "fyrtio", "-", "[>>]", 10},
			{50, "", "", "femtio", "-", "[>>]", 10},
			{60, "", "", "sextio", "-", "[>>]", 10},
			{70, "", "", "sjuttio", "-", "[>>]", 10},
			{80, "", "", "åttio", "-", "[>>]", 10},
			{90, "", "", "nittio", "-", "[>>]", 10},
			{100, "<<", "", "hundra", " ", "[>>]", 10},
			{1100, "<<", " ", "hundra", " ", "[>>]", 100},
			{2000, "<<", " ", "tusen", " ", "[>>]", 10},
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

	res, err = g.Expand(12, "default")
	exp = "tolv"
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Expand(1803, "default")
	exp = "arton hundra tre"
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Expand(1983, "default")
	exp = "nitton hundra åttio-tre"
	if res != exp {
		t.Errorf(fs, exp, res)
	}

	res, err = g.Expand(2001, "default")
	exp = "två tusen ett"
	if res != exp {
		t.Errorf(fs, exp, res)
	}

}
