package rbnf

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type RuleSet struct {
	Name  string
	Rules []BaseRule
}

type BaseRule struct {
	Base         int
	LeftSub      string
	LeftPadding  string
	SpellOut     string
	RightPadding string
	RightSub     string
	Radix        int
}

func (r BaseRule) Divisor() int {
	/** http://icu-project.org/apiref/icu4c/classRuleBasedNumberFormat.html
	To calculate the divisor, let [...] the exponent be the highest exponent of the radix that yields a result less than or equal to the base value.
	If the exponent is positive or 0, the divisor is the radix raised to the power of the exponent; otherwise, the divisor is 1.
	*/

	//for rad >= 0
	//exponent : the highest exponent of the radix that is less than or equal to the base value
	//divisor: radix^exponent
	var exponent, divisor int
	for i := 1; exp(r.Radix, i) <= r.Base; i++ {
		exponent = i
	}
	if exponent >= 0 {
		divisor = exp(r.Radix, exponent)
	} else {
		divisor = 1
	}
	return divisor
}
func (r BaseRule) Matches(input string) bool {
	n, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	return r.Base <= n
}

func (r BaseRule) Match(input string) MatchResult {
	n, err := strconv.Atoi(input)
	if err != nil {
		return MatchResult{}
	}
	divisor := r.Divisor()
	// >> in normal rule: Divide the number by the rule's divisor and format the remainder
	right := n % divisor
	// << in normal rule: Divide the number by the rule's divisor and format the quotient
	left := n / divisor
	return MatchResult{ForwardLeft: fmt.Sprintf("%d", left), ForwardRight: fmt.Sprintf("%d", right)}
}

type MatchResult struct {
	ForwardLeft  string
	ForwardRight string
}

type RulePackage struct {
	Language      string
	RuleSetGroups map[string]RuleSetGroup
}

func (p RulePackage) Spellout(input string, groupName string, ruleSetName string) (string, error) {
	if g, ok := p.RuleSetGroups[groupName]; ok {
		res, err := g.Spellout(input, ruleSetName)
		if err != nil {
			return "", err
		}
		return res, nil
	}
	return "", fmt.Errorf("No such rule set group: %s", groupName)
}

type RuleSetGroup struct {
	Name     string
	RuleSets map[string]RuleSet
}

func (g RuleSetGroup) FindRuleSet(ruleRef string) (RuleSet, bool) {
	if !isRuleRef(ruleRef) {
		return RuleSet{}, false
	}
	ruleName := strings.TrimPrefix(ruleRef, "%")
	res, ok := g.RuleSets[ruleName]
	return res, ok
}

func NewRuleSetGroup(name string, ruleSets []RuleSet) (RuleSetGroup, error) {
	rsMap := make(map[string]RuleSet)
	for _, rs := range ruleSets {
		// sort each rule set in ascending order
		sort.Slice(rs.Rules, func(i, j int) bool { return rs.Rules[i].Base < rs.Rules[j].Base })
		rsMap[rs.Name] = rs
	}
	res := RuleSetGroup{Name: name, RuleSets: rsMap}

	for _, ruleSet := range res.RuleSets {
		for _, rule := range ruleSet.Rules {
			if isRuleRef(rule.LeftSub) {
				if _, ok := res.FindRuleSet(rule.LeftSub); !ok {
					return res, fmt.Errorf("No such rule set: %s", rule.LeftSub)
				}
			}
			if isRuleRef(rule.RightSub) {
				if _, ok := res.FindRuleSet(rule.RightSub); !ok {
					return res, fmt.Errorf("No such rule set: %s", rule.RightSub)
				}
			}
		}
	}
	return res, nil
}

func findMatchingRule(input string, ruleSet RuleSet) (BaseRule, bool) {
	var res BaseRule
	var found = false
	for _, r := range ruleSet.Rules {
		if r.Matches(input) {
			res = r
			found = true
		} else {
			break
		}
	}
	return res, found
}

func (g RuleSetGroup) Spellout(input string, ruleSetName string) (string, error) {
	if rs, ok := g.RuleSets[ruleSetName]; ok {
		return g.spellout(input, rs), nil
	}
	return "", fmt.Errorf("No such rule set: %s", ruleSetName)
}

func (g RuleSetGroup) spellout(input string, ruleSet RuleSet) string {

	matchedRule, ok := findMatchingRule(input, ruleSet)
	if !ok {
		return input
	}

	if fmt.Sprintf("%d", matchedRule.Base) == input {
		// if n == 0 && matchedRule.Base == n {
		return matchedRule.SpellOut
	}

	match := matchedRule.Match(input)

	var left, right string
	if matchedRule.RightSub == "[>>]" { // Text in brackets is omitted if the number being formatted is an even multiple of 10
		//if n%10 != 0 {
		if !strings.HasSuffix(input, "0") {
			remSpelled := g.spellout(match.ForwardRight, ruleSet)
			right = remSpelled
		}
	} else if matchedRule.RightSub == ">>" {
		remSpelled := g.spellout(match.ForwardRight, ruleSet)
		right = remSpelled
	} else if namedRuleSet, ok := g.FindRuleSet((matchedRule.RightSub); ok {
		right = g.spellout(match.ForwardRight, namedRuleSet)
	} else if matchedRule.RightSub != "" {
		log.Fatalf("Unknown rule context: %s", matchedRule.RightSub)
	}

	if matchedRule.LeftSub == "<<" {
		left = g.spellout(match.ForwardLeft, ruleSet)
	} else if namedRuleSet, ok := g.FindRuleSet(matchedRule.LeftSub); ok {
		left = g.spellout(match.ForwardLeft, namedRuleSet)
	} else if matchedRule.LeftSub != "" {
		log.Fatalf("Unknown rule context: %s", matchedRule.LeftSub)
	}

	res := strings.TrimSpace(left + matchedRule.LeftPadding + matchedRule.SpellOut + matchedRule.RightPadding + right)
	return res
}

func exp(x, y int) int {
	res := 1
	for i := 1; i <= y; i++ {
		res = res * x
	}
	return res
}

func isRuleRef(s string) bool {
	return strings.HasPrefix(s, "%") // s != "" && !strings.Contains(s, "<") && !strings.Contains(s, ">")
}
