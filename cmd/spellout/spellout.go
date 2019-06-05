package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/stts-se/rbnf"
	"github.com/stts-se/rbnf/xmlreader"
)

func main() {
	cmd := path.Base(os.Args[0])

	// Flags
	var flags = flag.NewFlagSet(cmd, flag.ExitOnError)
	syntaxCheck := flags.Bool("s", false, "Check rule file syntax and exit")
	listRules := flags.Bool("l", false, "List rules and exit (rule groups and rule sets)")
	ruleGroup := flags.String("g", "SpelloutRules", "Use named `rule group`")
	ruleSet := flags.String("r", "spellout-numbering", "Use named `rule set`")
	help := flags.Bool("h", false, "Print usage and exit")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <options> <xml file/url> <input>\n", cmd)
		fmt.Fprintf(os.Stderr, "  if no input argument is specified, input will be read from stdin\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flags.PrintDefaults()
	}

	if *help || len(args) < 1 {
		flags.Usage()
		os.Exit(1)
	}

	var rPackage rbnf.RulePackage
	var err error
	f := args[0]
	if strings.HasPrefix(f, "http") {
		rPackage, err = xmlreader.RulesFromXMLURL(f)
	} else {
		rPackage, err = xmlreader.RulesFromXMLFile(f)
	}
	if err != nil {
		log.Fatalf("Couldn't parse file %s : %v", f, err)
	}
	log.Printf("Parsed xml rule file %s", f)

	foundRuleGroupAndRuleSet := false
	for _, g := range rPackage.RuleSetGroups {
		if g.Name == *ruleGroup {
			for _, s := range g.RuleSets {
				if s.Name == *ruleSet {
					foundRuleGroupAndRuleSet = true
				}
			}
		}
	}
	if !foundRuleGroupAndRuleSet {
		log.Fatalf("Couldn't find rule set %s/%s in rule file %s", *ruleGroup, *ruleSet, f)
	}

	if *syntaxCheck {
		os.Exit(0)
	}

	if *listRules {
		for _, g := range rPackage.RuleSetGroups {
			fmt.Printf("%s\n", g.Name)
			for _, s := range g.RuleSets {
				fmt.Printf(" - %s (%d)\n", s.Name, len(s.Rules))
			}
		}
		os.Exit(0)
	}

	var process = func(s string) {
		res, err := rPackage.Spellout(s, *ruleGroup, *ruleSet)
		if err != nil {
			log.Fatalf("Couldn't spellout %s : %v", s, err)
		}
		fmt.Printf("%s\t%s\n", s, res)
	}

	if len(args) == 1 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			process(s.Text())

		}
	} else { //if len(args) > 1 {
		for _, s := range args[1:] {
			process(s)
		}
	}

}
