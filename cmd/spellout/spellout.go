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

	spelloutCardinalReale := rbnf.RuleSet{
		Name: "spellout-cardinal-reale",
		Rules: []rbnf.BaseRule{
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

	spelloutCardinalNeuter := rbnf.RuleSet{
		Name: "spellout-cardinal-neuter",
		Rules: []rbnf.BaseRule{
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

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "USAGE: spellout <numerals>\n")
		os.Exit(1)
	}

	g, err := rbnf.NewRuleSetGroup(
		"spellout-cardinal",
		[]rbnf.RuleSet{
			defaultRules,
			spelloutCardinalReale,
			spelloutCardinalNeuter,
		})
	if err != nil {
		log.Fatalf("Couldn't create rule set group : %v", err)
	}

	for _, s := range os.Args[1:] {

		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Couldn't parse numeral %v : %v", s, err)
		}

		res, err := g.Spellout(n, "default")
		if err != nil {
			log.Fatalf("Couldn't spellout numeral %v : %v", n, err)
		}
		fmt.Printf("%s\n", res)
	}

}
