package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/stts-se/rbnf"
)

func main() {

	defaultRules := rbnf.RuleSet{
		Name: "default",
		Rules: []rbnf.BaseRule{
			rbnf.BaseRule{0, "", "", "noll", "", ""},
			rbnf.BaseRule{1, "", "", "ett", "", ""},
			rbnf.BaseRule{2, "", "", "två", "", ""},
			rbnf.BaseRule{3, "", "", "tre", "", ""},
			rbnf.BaseRule{4, "", "", "fyra", "", ""},
			rbnf.BaseRule{5, "", "", "fem", "", ""},
			rbnf.BaseRule{6, "", "", "sex", "", ""},
			rbnf.BaseRule{7, "", "", "sju", "", ""},
			rbnf.BaseRule{8, "", "", "åtta", "", ""},
			rbnf.BaseRule{9, "", "", "nio", "", ""},
			rbnf.BaseRule{10, "", "", "tio", "", ""},
			rbnf.BaseRule{11, "", "", "elva", "", ""},
			rbnf.BaseRule{12, "", "", "tolv", "", ""},
			rbnf.BaseRule{13, "", "", "tretton", "", ""},
			rbnf.BaseRule{14, "", "", "fjorton", "", ""},
			rbnf.BaseRule{15, "", "", "femton", "", ""},
			rbnf.BaseRule{16, "", "", "sexton", "", ""},
			rbnf.BaseRule{17, "", "", "sjutton", "", ""},
			rbnf.BaseRule{18, "", "", "arton", "", ""},
			rbnf.BaseRule{19, "", "", "nitton", "", ""},
			rbnf.BaseRule{20, "", "", "tjugo", "-", "[>>]"},
			rbnf.BaseRule{30, "", "", "trettio", "-", "[>>]"},
			rbnf.BaseRule{40, "", "", "fyrtio", "-", "[>>]"},
			rbnf.BaseRule{50, "", "", "femtio", "-", "[>>]"},
			rbnf.BaseRule{60, "", "", "sextio", "-", "[>>]"},
			rbnf.BaseRule{70, "", "", "sjuttio", "-", "[>>]"},
			rbnf.BaseRule{80, "", "", "åttio", "-", "[>>]"},
			rbnf.BaseRule{90, "", "", "nittio", "-", "[>>]"},

			rbnf.BaseRule{100, "<<", " ", "hundra", " ", "[>>]"},

			rbnf.BaseRule{1000, "", " ", "ettusen", " ", "[>>]"},
			rbnf.BaseRule{2000, "spellout-cardinal-reale", " ", "tusen", " ", "[>>]"},

			rbnf.BaseRule{1000000, "", " ", "en miljon", " ", "[>>]"},
			rbnf.BaseRule{2000000, "spellout-cardinal-reale", " ", "miljoner", " ", "[>>]"},
			rbnf.BaseRule{1000000000, "", "", "en miljard", " ", "[>>]"},
			rbnf.BaseRule{2000000000, "spellout-cardinal-reale", " ", "miljarder", " ", "[>>]"},
		},
	}

	spelloutCardinalReale := rbnf.RuleSet{
		Name: "spellout-cardinal-reale",
		Rules: []rbnf.BaseRule{
			rbnf.BaseRule{0, "", "", "noll", "", ""},
			rbnf.BaseRule{1, "", "", "en", "", ""},
			rbnf.BaseRule{2, "", "", "två", "", ""},
			rbnf.BaseRule{3, "", "", "tre", "", ""},
			rbnf.BaseRule{4, "", "", "fyra", "", ""},
			rbnf.BaseRule{5, "", "", "fem", "", ""},
			rbnf.BaseRule{6, "", "", "sex", "", ""},
			rbnf.BaseRule{7, "", "", "sju", "", ""},
			rbnf.BaseRule{8, "", "", "åtta", "", ""},
			rbnf.BaseRule{9, "", "", "nio", "", ""},
			rbnf.BaseRule{10, "", "", "tio", "", ""},
			rbnf.BaseRule{11, "", "", "elva", "", ""},
			rbnf.BaseRule{12, "", "", "tolv", "", ""},
			rbnf.BaseRule{13, "", "", "tretton", "", ""},
			rbnf.BaseRule{14, "", "", "fjorton", "", ""},
			rbnf.BaseRule{15, "", "", "femton", "", ""},
			rbnf.BaseRule{16, "", "", "sexton", "", ""},
			rbnf.BaseRule{17, "", "", "sjutton", "", ""},
			rbnf.BaseRule{18, "", "", "arton", "", ""},
			rbnf.BaseRule{19, "", "", "nitton", "", ""},
			rbnf.BaseRule{20, "", "", "tjugo", "-", "[>>]"},
			rbnf.BaseRule{30, "", "", "trettio", "-", "[>>]"},
			rbnf.BaseRule{40, "", "", "fyrtio", "-", "[>>]"},
			rbnf.BaseRule{50, "", "", "femtio", "-", "[>>]"},
			rbnf.BaseRule{60, "", "", "sextio", "-", "[>>]"},
			rbnf.BaseRule{70, "", "", "sjuttio", "-", "[>>]"},
			rbnf.BaseRule{80, "", "", "åttio", "-", "[>>]"},
			rbnf.BaseRule{90, "", "", "nittio", "-", "[>>]"},
			rbnf.BaseRule{100, "spellout-cardinal-neuter", " ", "hundra", " ", "[>>]"},
			rbnf.BaseRule{1000, "", " ", "ettusen", "-", "[>>]"},
			rbnf.BaseRule{2000, "spellout-cardinal-reale", " ", "tusen", " ", "[>>]"},
			rbnf.BaseRule{1000000, "", " ", "en miljon", " ", "[>>]"},
			rbnf.BaseRule{2000000, "spellout-cardinal-reale", " ", "miljoner", " ", "[>>]"},
		},
	}

	spelloutCardinalNeuter := rbnf.RuleSet{
		Name: "spellout-cardinal-neuter",
		Rules: []rbnf.BaseRule{
			rbnf.BaseRule{0, "", "", "noll", "", ""},
			rbnf.BaseRule{1, "", "", "ett", "", ""},
			rbnf.BaseRule{2, "", "", "två", "", ""},
			rbnf.BaseRule{3, "", "", "tre", "", ""},
			rbnf.BaseRule{4, "", "", "fyra", "", ""},
			rbnf.BaseRule{5, "", "", "fem", "", ""},
			rbnf.BaseRule{6, "", "", "sex", "", ""},
			rbnf.BaseRule{7, "", "", "sju", "", ""},
			rbnf.BaseRule{8, "", "", "åtta", "", ""},
			rbnf.BaseRule{9, "", "", "nio", "", ""},
			rbnf.BaseRule{10, "", "", "tio", "", ""},
			rbnf.BaseRule{11, "", "", "elva", "", ""},
			rbnf.BaseRule{12, "", "", "tolv", "", ""},
			rbnf.BaseRule{13, "", "", "tretton", "", ""},
			rbnf.BaseRule{14, "", "", "fjorton", "", ""},
			rbnf.BaseRule{15, "", "", "femton", "", ""},
			rbnf.BaseRule{16, "", "", "sexton", "", ""},
			rbnf.BaseRule{17, "", "", "sjutton", "", ""},
			rbnf.BaseRule{18, "", "", "arton", "", ""},
			rbnf.BaseRule{19, "", "", "nitton", "", ""},
			rbnf.BaseRule{20, "", "", "tjugo", "-", "[>>]"},
			rbnf.BaseRule{30, "", "", "trettio", "-", "[>>]"},
			rbnf.BaseRule{40, "", "", "fyrtio", "-", "[>>]"},
			rbnf.BaseRule{50, "", "", "femtio", "-", "[>>]"},
			rbnf.BaseRule{60, "", "", "sextio", "-", "[>>]"},
			rbnf.BaseRule{70, "", "", "sjuttio", "-", "[>>]"},
			rbnf.BaseRule{80, "", "", "åttio", "-", "[>>]"},
			rbnf.BaseRule{90, "", "", "nittio", "-", "[>>]"},
			rbnf.BaseRule{100, "spellout-cardinal-neuter", "", "hundra", " ", "[>>]"},
			rbnf.BaseRule{1000, "", "", "ettusen", " ", "[>>]"},
			rbnf.BaseRule{2000, "spellout-cardinal-reale", "", "tusen", " ", "[>>]"},
			rbnf.BaseRule{1000000, "", "", "en miljon", " ", "[>>]"},
			rbnf.BaseRule{2000000, "spellout-cardinal-reale", " ", "miljoner", " ", "[>>]"},
			rbnf.BaseRule{1000000000, "", "", "en miljard", " ", "[>>]"},
			rbnf.BaseRule{2000000000, "spellout-cardinal-reale", " ", "miljarder", " ", "[>>]"},
		},
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "USAGE: spellout <numerals>\n")
		os.Exit(1)
	}

	for _, s := range os.Args[1:] {

		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Couldn't parse numeral %v : %v", s, err)
		}

		g, err := rbnf.NewRuleSetGroup(
			"spellout-cardinal",
			map[string]rbnf.RuleSet{
				defaultRules.Name:           defaultRules,
				spelloutCardinalReale.Name:  spelloutCardinalReale,
				spelloutCardinalNeuter.Name: spelloutCardinalNeuter,
			})
		if err != nil {
			log.Fatalf("Couldn't create rule set group : %v", err)
		}

		res, err := g.Expand(n, "default")
		if err != nil {
			log.Fatalf("Couldn't expand numeral %v : %v", n, err)
		}
		fmt.Printf("%s\n", res)
	}

}
