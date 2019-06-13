package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/stts-se/rbnf"
	"github.com/stts-se/rbnf/xmlreader"
)

func main() {
	cmd := path.Base(os.Args[0])

	// Flags
	var flags = flag.NewFlagSet(cmd, flag.ExitOnError)
	syntaxCheck := flags.Bool("s", false, "Check rule file syntax and exit")
	listPublicRules := flags.Bool("l", false, "List public rules and exit (rule groups and rule sets)")
	listAllRules := flags.Bool("L", false, "List all (private/public) rules and exit (rule groups and rule sets)")
	ruleGroup := flags.String("g", "", "Use named `rule group` (default first group)")
	ruleSet := flags.String("r", "", "Use named `rule set`")
	debug := flags.Bool("d", false, "Debug")
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
		os.Exit(2)
	}

	var rPackage rbnf.RulePackage
	var err error
	f := args[0]

	xmlreader.Verb = *debug

	if strings.HasPrefix(f, "http") {
		rPackage, err = xmlreader.RulesFromXMLURL(f)
	} else {
		rPackage, err = xmlreader.RulesFromXMLFile(f)
	}
	rPackage.Debug = *debug
	if err != nil {
		log.Fatalf("Couldn't parse file %s : %v", f, err)
	}
	log.Printf("Parsed rule file %s", f)

	if *syntaxCheck && *ruleSet == "" {
		os.Exit(0)
	}

	if *listAllRules || *listPublicRules {
		if *listAllRules {
			fmt.Println("== Listing all rule sets ==")
		} else {
			fmt.Println("== Listing public rule sets ==")
		}
		for _, g := range rPackage.RuleSetGroups {
			fmt.Printf("%s\n", g.Name)
			rs := []rbnf.RuleSet{}
			for _, s := range g.RuleSets {
				if *listAllRules || !s.Private {
					rs = append(rs, s)
				}
			}
			sort.Slice(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name })
			for _, s := range rs {
				access := "public"
				if s.Private {
					access = "private"
				}
				fmt.Printf(" - %s [%s] (%d)\n", s.Name, access, len(s.Rules))
			}
		}
		os.Exit(0)
	}

	// validate specified rule group and rule set
	if *ruleSet == "" {
		fmt.Fprintf(os.Stderr, "flag -r (rule set) is required\n")
		flags.Usage()
		os.Exit(1)
	}
	foundRuleGroupAndRuleSet := false
	if *ruleGroup == "" {
		ruleGroup = &rPackage.RuleSetGroups[0].Name
	}
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

	var nSpelled = 0
	var process = func(s string) {
		res, err := rPackage.Spellout(s, *ruleGroup, *ruleSet, *debug)
		nSpelled++
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

	///fmt.Fprintf(os.Stderr, "[%s] No of spelled numerals: %v\n", cmd, nSpelled)

}
