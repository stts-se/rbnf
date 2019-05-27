package rbnf

import (
	"fmt"
	"log"
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
}

type RuleSetGroup struct {
	Name     string
	RuleSets map[string]RuleSet
}

func NewRuleSetGroup(name string, ruleSets map[string]RuleSet) (RuleSetGroup, error) {
	for _, ruleSet := range ruleSets {
		for _, rule := range ruleSet.Rules {
			if isRuleName(rule.LeftPadding) {
			}
		}
	}
	return RuleSetGroup{Name: name, RuleSets: ruleSets}, nil
}

func (r BaseRule) Divisor() int {
	/** http://icu-project.org/apiref/icu4c/classRuleBasedNumberFormat.html
	bv/rad>: 	bv specifies the rule's base value.

	To calculate the divisor, let the radix be rad, and the exponent be the highest exponent of the radix that yields a result less than or equal to the base value.
	If the exponent is positive or 0, the divisor is the radix raised to the power of the exponent; otherwise, the divisor is 1.
	*/

	radix := 10 // TODO: variable radix

	//for rad >= 0
	//exponent : the highest exponent of the radix that is less than or equal to the base value
	//divisor: radix^exponent
	var exponent, divisor int
	for i := 1; exp(radix, i) <= r.BaseNum; i++ {
		exponent = i
	}
	if exponent >= 0 {
		divisor = exp(radix, exponent)
	} else {
		divisor = 1
	}
	return divisor
}

func (g RuleSetGroup) Expand(n int, ruleSet string) (string, error) {
	if rs, ok := g.RuleSets[ruleSet]; ok {
		return g.expand(n, rs), nil
	}
	return "", fmt.Errorf("No such rule set: %s", ruleSet)
}

func (g RuleSetGroup) expand(n int, ruleSet RuleSet) string {

	// TODO: sort rules in ascending order
	var factor BaseRule
	for _, r := range ruleSet.Rules {
		if r.BaseNum <= n {
			factor = r
		}
		if r.BaseNum > n {
			break
		}
	}

	if n == 0 && factor.BaseNum == n {
		return factor.SpellOut
	}

	divisor := factor.Divisor()

	// >> in normal rule: Divide the number by the rule's divisor and format the remainder
	remainderRight := n % divisor

	// << in normal rule: Divide the number by the rule's divisor and format the quotient
	remainderLeft := n / divisor

	var left, right string
	if remainderRight != 0 {
		if factor.RightSub == "[>>]" { // Text in brackets is omitted if the number being formatted is an even multiple of 10
			if n%10 != 0 {
				remSpelled := g.expand(remainderRight, ruleSet)
				right = remSpelled
			}
		} else if factor.RightSub == ">>" {
			remSpelled := g.expand(remainderRight, ruleSet)
			right = remSpelled
		} else if namedRuleSet, ok := g.RuleSets[factor.RightSub]; ok {
			right = g.expand(remainderRight, namedRuleSet)
		} else if factor.RightSub != "" {
			log.Fatalf("Unknown rule context: %s", factor.RightSub)
		}
	}

	if remainderLeft != 0 {
		if factor.LeftSub == "<<" {
			left = g.expand(remainderLeft, ruleSet)
		} else if namedRuleSet, ok := g.RuleSets[factor.LeftSub]; ok {
			left = g.expand(remainderLeft, namedRuleSet)
		} else if factor.LeftSub != "" {
			log.Fatalf("Unknown rule context: %s", factor.LeftSub)
		}
	}

	res := strings.TrimSpace(left + factor.LeftPadding + factor.SpellOut + factor.RightPadding + right)
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
