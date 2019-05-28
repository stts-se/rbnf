package rbnf

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type RuleSet struct {
	Name  string
	Rules []BaseRule
}

type BaseRule struct {
	BaseNum      int
	LeftSub      string
	LeftPadding  string
	SpellOut     string
	RightPadding string
	RightSub     string
	Radix        int
}

// type IntRule struct {
// 	BaseNum      int
// 	LeftSub      string
// 	LeftPadding  string
// 	SpellOut     string
// 	RightPadding string
// 	RightSub     string
// 	Radix        int
// }

// type MatchRes struct {
// 	Left   string
// 	Middle string
// 	Right  string
// }

// type BaseRule interface {
// 	Match(input string) MatchRes
// }

type RuleSetGroup struct {
	Name     string
	RuleSets map[string]RuleSet
}

func NewRuleSetGroup(name string, ruleSets []RuleSet) (RuleSetGroup, error) {
	rsMap := make(map[string]RuleSet)
	for _, rs := range ruleSets {
		// sort each rule set in ascending order
		sort.Slice(rs.Rules, func(i, j int) bool { return rs.Rules[i].BaseNum < rs.Rules[j].BaseNum })
		rsMap[rs.Name] = rs
	}
	res := RuleSetGroup{Name: name, RuleSets: rsMap}

	for _, ruleSet := range res.RuleSets {
		for _, rule := range ruleSet.Rules {
			if isRuleName(rule.LeftSub) {
				if _, ok := res.RuleSets[rule.LeftSub]; !ok {
					return res, fmt.Errorf("No such rule set: %s", rule.LeftSub)
				}
			}
			if isRuleName(rule.RightSub) {
				if _, ok := res.RuleSets[rule.RightSub]; !ok {
					return res, fmt.Errorf("No such rule set: %s", rule.RightSub)
				}
			}
		}
	}
	return res, nil
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
	for i := 1; exp(r.Radix, i) <= r.BaseNum; i++ {
		exponent = i
	}
	if exponent >= 0 {
		divisor = exp(r.Radix, exponent)
	} else {
		divisor = 1
	}
	return divisor
}

func (g RuleSetGroup) Spellout(n int, ruleSet string) (string, error) {
	if rs, ok := g.RuleSets[ruleSet]; ok {
		return g.spellout(n, rs), nil
	}
	return "", fmt.Errorf("No such rule set: %s", ruleSet)
}

func findMatchingRule(n int, ruleSet RuleSet) (BaseRule, bool) {
	var res BaseRule
	var found = false
	for _, r := range ruleSet.Rules {
		if r.BaseNum <= n {
			res = r
			found = true
		}
		if r.BaseNum > n {
			break
		}
	}
	return res, found
}

func (g RuleSetGroup) spellout(n int, ruleSet RuleSet) string {

	matchedRule, ok := findMatchingRule(n, ruleSet)
	if !ok {
		return fmt.Sprintf("%d", n)
	}

	if n == 0 && matchedRule.BaseNum == n {
		return matchedRule.SpellOut
	}

	divisor := matchedRule.Divisor()

	// >> in normal rule: Divide the number by the rule's divisor and format the remainder
	remainderRight := n % divisor

	// << in normal rule: Divide the number by the rule's divisor and format the quotient
	remainderLeft := n / divisor

	var left, right string
	if remainderRight != 0 {
		if matchedRule.RightSub == "[>>]" { // Text in brackets is omitted if the number being formatted is an even multiple of 10
			if n%10 != 0 {
				remSpelled := g.spellout(remainderRight, ruleSet)
				right = remSpelled
			}
		} else if matchedRule.RightSub == ">>" {
			remSpelled := g.spellout(remainderRight, ruleSet)
			right = remSpelled
		} else if namedRuleSet, ok := g.RuleSets[matchedRule.RightSub]; ok {
			right = g.spellout(remainderRight, namedRuleSet)
		} else if matchedRule.RightSub != "" {
			log.Fatalf("Unknown rule context: %s", matchedRule.RightSub)
		}
	}

	if remainderLeft != 0 {
		if matchedRule.LeftSub == "<<" {
			left = g.spellout(remainderLeft, ruleSet)
		} else if namedRuleSet, ok := g.RuleSets[matchedRule.LeftSub]; ok {
			left = g.spellout(remainderLeft, namedRuleSet)
		} else if matchedRule.LeftSub != "" {
			log.Fatalf("Unknown rule context: %s", matchedRule.LeftSub)
		}
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

func isRuleName(s string) bool {
	return s != "" && !strings.Contains(s, "<") && !strings.Contains(s, ">")
}
