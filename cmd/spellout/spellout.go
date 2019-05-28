package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stts-se/rbnf"
)

func main() {

	defaultRules := rbnf.RuleSet{
		Name: "default",
		Rules: []rbnf.BaseRule{
			rbnf.NewStringRule("-x", "", "", "minus", " ", ">>"),
			rbnf.NewStringRule("x.x", "<<", " ", "komma", " ", ">>"),
			rbnf.NewStringRule("x,x", "<<", " ", "komma", " ", ">>"),
			rbnf.NewStringRule("x%", "<<", " ", "procent", " ", ""),
			rbnf.NewStringRule("x‰", "<<", " ", "promille", " ", ""),
			rbnf.NewIntRule(0, "", "", "noll", "", "", 10),
			rbnf.NewIntRule(1, "", "", "ett", "", "", 10),
			rbnf.NewIntRule(2, "", "", "två", "", "", 10),
			rbnf.NewIntRule(3, "", "", "tre", "", "", 10),
			rbnf.NewIntRule(4, "", "", "fyra", "", "", 10),
			rbnf.NewIntRule(5, "", "", "fem", "", "", 10),
			rbnf.NewIntRule(6, "", "", "sex", "", "", 10),
			rbnf.NewIntRule(7, "", "", "sju", "", "", 10),
			rbnf.NewIntRule(8, "", "", "åtta", "", "", 10),
			rbnf.NewIntRule(9, "", "", "nio", "", "", 10),
			rbnf.NewIntRule(10, "", "", "tio", "", "", 10),
			rbnf.NewIntRule(11, "", "", "elva", "", "", 10),
			rbnf.NewIntRule(12, "", "", "tolv", "", "", 10),
			rbnf.NewIntRule(13, "", "", "tretton", "", "", 10),
			rbnf.NewIntRule(14, "", "", "fjorton", "", "", 10),
			rbnf.NewIntRule(15, "", "", "femton", "", "", 10),
			rbnf.NewIntRule(16, "", "", "sexton", "", "", 10),
			rbnf.NewIntRule(17, "", "", "sjutton", "", "", 10),
			rbnf.NewIntRule(18, "", "", "arton", "", "", 10),
			rbnf.NewIntRule(19, "", "", "nitton", "", "", 10),
			rbnf.NewIntRule(20, "", "", "tjugo", "-", "[>>]", 10),
			rbnf.NewIntRule(30, "", "", "trettio", "-", "[>>]", 10),
			rbnf.NewIntRule(40, "", "", "fyrtio", "-", "[>>]", 10),
			rbnf.NewIntRule(50, "", "", "femtio", "-", "[>>]", 10),
			rbnf.NewIntRule(60, "", "", "sextio", "-", "[>>]", 10),
			rbnf.NewIntRule(70, "", "", "sjuttio", "-", "[>>]", 10),
			rbnf.NewIntRule(80, "", "", "åttio", "-", "[>>]", 10),
			rbnf.NewIntRule(90, "", "", "nittio", "-", "[>>]", 10),

			rbnf.NewIntRule(100, "<<", " ", "hundra", " ", "[>>]", 10),

			rbnf.NewIntRule(1000, "", " ", "ettusen", " ", "[>>]", 10),
			rbnf.NewIntRule(2000, "%spellout-cardinal-reale", " ", "tusen", " ", "[>>]", 10),

			rbnf.NewIntRule(1000000, "", " ", "en miljon", " ", "[>>]", 10),
			rbnf.NewIntRule(2000000, "%spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10),
			rbnf.NewIntRule(1000000000, "", "", "en miljard", " ", "[>>]", 10),
			rbnf.NewIntRule(2000000000, "%spellout-cardinal-reale", " ", "miljarder", " ", "[>>]", 10),
		},
	}

	spelloutCardinalReale := rbnf.RuleSet{
		Name: "spellout-cardinal-reale",
		Rules: []rbnf.BaseRule{
			rbnf.NewIntRule(0, "", "", "noll", "", "", 10),
			rbnf.NewIntRule(1, "", "", "en", "", "", 10),
			rbnf.NewIntRule(2, "", "", "två", "", "", 10),
			rbnf.NewIntRule(3, "", "", "tre", "", "", 10),
			rbnf.NewIntRule(4, "", "", "fyra", "", "", 10),
			rbnf.NewIntRule(5, "", "", "fem", "", "", 10),
			rbnf.NewIntRule(6, "", "", "sex", "", "", 10),
			rbnf.NewIntRule(7, "", "", "sju", "", "", 10),
			rbnf.NewIntRule(8, "", "", "åtta", "", "", 10),
			rbnf.NewIntRule(9, "", "", "nio", "", "", 10),
			rbnf.NewIntRule(10, "", "", "tio", "", "", 10),
			rbnf.NewIntRule(11, "", "", "elva", "", "", 10),
			rbnf.NewIntRule(12, "", "", "tolv", "", "", 10),
			rbnf.NewIntRule(13, "", "", "tretton", "", "", 10),
			rbnf.NewIntRule(14, "", "", "fjorton", "", "", 10),
			rbnf.NewIntRule(15, "", "", "femton", "", "", 10),
			rbnf.NewIntRule(16, "", "", "sexton", "", "", 10),
			rbnf.NewIntRule(17, "", "", "sjutton", "", "", 10),
			rbnf.NewIntRule(18, "", "", "arton", "", "", 10),
			rbnf.NewIntRule(19, "", "", "nitton", "", "", 10),
			rbnf.NewIntRule(20, "", "", "tjugo", "-", "[>>]", 10),
			rbnf.NewIntRule(30, "", "", "trettio", "-", "[>>]", 10),
			rbnf.NewIntRule(40, "", "", "fyrtio", "-", "[>>]", 10),
			rbnf.NewIntRule(50, "", "", "femtio", "-", "[>>]", 10),
			rbnf.NewIntRule(60, "", "", "sextio", "-", "[>>]", 10),
			rbnf.NewIntRule(70, "", "", "sjuttio", "-", "[>>]", 10),
			rbnf.NewIntRule(80, "", "", "åttio", "-", "[>>]", 10),
			rbnf.NewIntRule(90, "", "", "nittio", "-", "[>>]", 10),
			rbnf.NewIntRule(100, "%spellout-cardinal-neuter", " ", "hundra", " ", "[>>]", 10),
			rbnf.NewIntRule(1000, "", " ", "ettusen", "-", "[>>]", 10),
			rbnf.NewIntRule(2000, "%spellout-cardinal-reale", " ", "tusen", " ", "[>>]", 10),
			rbnf.NewIntRule(1000000, "", " ", "en miljon", " ", "[>>]", 10),
			rbnf.NewIntRule(2000000, "%spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10),
		},
	}

	spelloutCardinalNeuter := rbnf.RuleSet{
		Name: "spellout-cardinal-neuter",
		Rules: []rbnf.BaseRule{
			rbnf.NewIntRule(0, "", "", "noll", "", "", 10),
			rbnf.NewIntRule(1, "", "", "ett", "", "", 10),
			rbnf.NewIntRule(2, "", "", "två", "", "", 10),
			rbnf.NewIntRule(3, "", "", "tre", "", "", 10),
			rbnf.NewIntRule(4, "", "", "fyra", "", "", 10),
			rbnf.NewIntRule(5, "", "", "fem", "", "", 10),
			rbnf.NewIntRule(6, "", "", "sex", "", "", 10),
			rbnf.NewIntRule(7, "", "", "sju", "", "", 10),
			rbnf.NewIntRule(8, "", "", "åtta", "", "", 10),
			rbnf.NewIntRule(9, "", "", "nio", "", "", 10),
			rbnf.NewIntRule(10, "", "", "tio", "", "", 10),
			rbnf.NewIntRule(11, "", "", "elva", "", "", 10),
			rbnf.NewIntRule(12, "", "", "tolv", "", "", 10),
			rbnf.NewIntRule(13, "", "", "tretton", "", "", 10),
			rbnf.NewIntRule(14, "", "", "fjorton", "", "", 10),
			rbnf.NewIntRule(15, "", "", "femton", "", "", 10),
			rbnf.NewIntRule(16, "", "", "sexton", "", "", 10),
			rbnf.NewIntRule(17, "", "", "sjutton", "", "", 10),
			rbnf.NewIntRule(18, "", "", "arton", "", "", 10),
			rbnf.NewIntRule(19, "", "", "nitton", "", "", 10),
			rbnf.NewIntRule(20, "", "", "tjugo", "-", "[>>]", 10),
			rbnf.NewIntRule(30, "", "", "trettio", "-", "[>>]", 10),
			rbnf.NewIntRule(40, "", "", "fyrtio", "-", "[>>]", 10),
			rbnf.NewIntRule(50, "", "", "femtio", "-", "[>>]", 10),
			rbnf.NewIntRule(60, "", "", "sextio", "-", "[>>]", 10),
			rbnf.NewIntRule(70, "", "", "sjuttio", "-", "[>>]", 10),
			rbnf.NewIntRule(80, "", "", "åttio", "-", "[>>]", 10),
			rbnf.NewIntRule(90, "", "", "nittio", "-", "[>>]", 10),
			rbnf.NewIntRule(100, "%spellout-cardinal-neuter", "", "hundra", " ", "[>>]", 10),
			rbnf.NewIntRule(1000, "", "", "ettusen", " ", "[>>]", 10),
			rbnf.NewIntRule(2000, "%spellout-cardinal-reale", "", "tusen", " ", "[>>]", 10),
			rbnf.NewIntRule(1000000, "", "", "en miljon", " ", "[>>]", 10),
			rbnf.NewIntRule(2000000, "%spellout-cardinal-reale", " ", "miljoner", " ", "[>>]", 10),
			rbnf.NewIntRule(1000000000, "", "", "en miljard", " ", "[>>]", 10),
			rbnf.NewIntRule(2000000000, "%spellout-cardinal-reale", " ", "miljarder", " ", "[>>]", 10),
		},
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "USAGE: spellout <numerals>\n")
		os.Exit(1)
	}

	g, err := rbnf.NewRuleSetGroup(
		"SpelloutRules",
		[]rbnf.RuleSet{
			defaultRules,
			spelloutCardinalReale,
			spelloutCardinalNeuter,
		})
	if err != nil {
		log.Fatalf("Couldn't create rule set group : %v", err)
	}

	for _, s := range os.Args[1:] {

		res, err := g.Spellout(s, "default")
		if err != nil {
			log.Fatalf("Couldn't spellout numeral %v : %v", s, err)
		}
		fmt.Printf("%s\n", res)
	}

}
