package xmlreader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/stts-se/rbnf"
)

func readXMLFile(fn string) (Ldml, error) {
	res := Ldml{}

	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return res, fmt.Errorf("failed to read XML file : %v", err)
	}

	err = xml.Unmarshal(bytes, &res)
	if err != nil {
		return res, fmt.Errorf("failed to peocess XML file : %v", err)
	}

	return res, nil
}

func convertRuleSet(rs *Ruleset) (rbnf.RuleSet, error) {
	var res rbnf.RuleSet
	res.Name = rs.Attrtype
	for _, r := range rs.Rbnfrule {
		//fmt.Printf("RULE %#v\n", r)
		rule := rbnf.BaseRule{}
		baseNum, err := strconv.Atoi(r.Attrvalue)
		if err == nil { // numeric rule
			rule.BaseInt = baseNum
			// TODO test
			if r.Attrradix != "" {
				radix, err := strconv.Atoi(r.Attrradix)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to convert radix : %v\n", err)
				} else {
					rule.Radix = radix
				}
			} else {
				rule.Radix = 10 // Default radix
			}
		} else { // non-numeric rule
			rule.BaseString = r.Attrvalue
		}
		// TODO parse string
		// examples:
		// två;
		// trettio[­→→];
		// ←%spellout-cardinal-reale← miljoner[ →→];
		// minus →→;
		// ←← komma →→;

		//rule.LeftSub = "ls"   //r.String
		//rule.RightSub = "rs"  //r.String
		//rule.SpellOut = "apa" //r.String

		fmt.Println(rule)
		res.Rules = append(res.Rules, rule)
	}

	return res, nil
}

func convertGroup(g *RulesetGrouping) (string, []rbnf.RuleSet, error) {
	var res []rbnf.RuleSet
	name := g.Attrtype
	if strings.TrimSpace(name) == "" {
		return "", res, fmt.Errorf("rule set grouping lacks type attribute value")
	}

	for _, rs := range g.Ruleset {
		rbntRuleSet, err := convertRuleSet(rs)
		if err != nil {
			return rbntRuleSet.Name, res, fmt.Errorf("failed to convert rule set : %v", err)
			//fmt.Fprintf(os.Stderr, "skipping rule set '%s' : %v\n", rbntRuleSet.Name, err)
			//continue
		}
		res = append(res, rbntRuleSet)
	}

	return name, res, nil
}

func rulesFromLdml(ldml Ldml) ([]rbnf.RuleSetGroup, error) {
	//var res rbnf.RuleSetGroup
	res := []rbnf.RuleSetGroup{}
	// name of whole file

	groups := ldml.Rbnf.RulesetGrouping
	if len(groups) == 0 {
		return res, fmt.Errorf("empty RulsetGrouping")
	}

	var rbnfGroups []rbnf.RuleSetGroup
	for _, g := range groups {
		name, ruleSet, err := convertGroup(g)
		if err != nil {
			return res, fmt.Errorf("failed to convert rule group : %v", err)
			//fmt.Fprintf(os.Stderr, "skipping rule group '%s' : %v", name, err)
			//continue
		}
		group, err := rbnf.NewRuleSetGroup(name, ruleSet)
		if err != nil {

			fmt.Printf("%#v\n", group)
			return res, fmt.Errorf("failed creating rbnf.NewRuleSetGroup instance : %v", err)

			//fmt.Fprintf(os.Stderr, "skipping rules set group '%s' : %v\n", name, err)
			//continue
		}
		rbnfGroups = append(rbnfGroups, group)
	}

	for _, g := range rbnfGroups {
		res = append(res, g)
	}

	return res, nil
}

func RulesFromXMLFile(fn string) (string, []rbnf.RuleSetGroup, error) {
	var lang string
	var res []rbnf.RuleSetGroup

	ldml, err := readXMLFile(fn)
	if err != nil {
		return lang, res, fmt.Errorf("RulesFromXMLFile: %v", err)
	}
	lang = ldml.Identity.Language.Attrtype

	res, err = rulesFromLdml(ldml)
	if err != nil {
		return lang, res, fmt.Errorf("RulesFromXMLFile: %v", err)
	}

	return lang, res, nil
}
