package xmlreader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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
	for _, r := range rs.Rbnfrule {
		fmt.Printf("RULE %#v\n", r)
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
			return name, res, fmt.Errorf("failed to convert rule set : %v", err)
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
		}
		group, err := rbnf.NewRuleSetGroup(name, ruleSet)
		if err != nil {
			return res, fmt.Errorf("failed creating rbnf.NewRuleSetGroup instance : %v", err)
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
