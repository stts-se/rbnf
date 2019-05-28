package rbnf

import (
	"fmt"
	"regexp"
	// "sort"
	"strconv"
	"strings"
)

type RuleSet struct {
	Name  string
	Rules []BaseRule
}

type BaseRule struct {
	BaseInt      int
	BaseString   string
	LeftSub      string
	LeftPadding  string
	SpellOut     string
	RightPadding string
	RightSub     string
	Radix        int

	stringMatchRegexp *regexp.Regexp
}

func NewIntRule(baseInt int, leftSub, leftPadding, spellOut, rightPadding, rightSub string, radix int) BaseRule {
	return BaseRule{
		BaseInt:      baseInt,
		LeftSub:      leftSub,
		LeftPadding:  leftPadding,
		SpellOut:     spellOut,
		RightPadding: rightPadding,
		RightSub:     rightSub,
		Radix:        radix,
	}
}
func NewStringRule(baseString string, leftSub, leftPadding, spellOut, rightPadding, rightSub string) BaseRule {
	return BaseRule{
		BaseString:        baseString,
		LeftSub:           leftSub,
		LeftPadding:       leftPadding,
		SpellOut:          spellOut,
		RightPadding:      rightPadding,
		RightSub:          rightSub,
		stringMatchRegexp: buildStringMatchRegexp(baseString),
	}
}

func (r BaseRule) String() string {
	if r.IsIntRule() {
		return fmt.Sprintf("%d => '%s%s' <%s> '%s%s' [r:%d]", r.BaseInt, r.LeftSub, r.LeftPadding, r.SpellOut, r.RightPadding, r.RightSub, r.Radix)
	}
	return fmt.Sprintf("%s => '%s%s' <%s> '%s%s'", r.BaseString, r.LeftSub, r.LeftPadding, r.SpellOut, r.RightPadding, r.RightSub)
}

func (r BaseRule) IsIntRule() bool {
	return r.BaseString == ""
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
	for i := 1; exp(r.Radix, i) <= r.BaseInt; i++ {
		exponent = i
	}
	if exponent >= 0 {
		divisor = exp(r.Radix, exponent)
	} else {
		divisor = 1
	}
	return divisor
}

func regexpEscape(s string) string {
	res := s
	chars := []string{`]`, `^`, `\`, `[`, `.`, `(`, `)`, `-`}
	for _, ch := range chars {
		res = strings.ReplaceAll(res, ch, fmt.Sprintf(`\%s`, ch))
	}
	return res
}

var nonXRE = regexp.MustCompile("([^x]+)")
var noInitialX = regexp.MustCompile("^([^x])")
var noFinalX = regexp.MustCompile("([^x])$")
var emptyRegexp *regexp.Regexp

// TODO: this is sooo ugly -- can it be done better?
func buildStringMatchRegexp(baseString string) *regexp.Regexp {
	reString := baseString
	reString = regexpEscape(reString)                        // escape special chars in the BaseString
	reString = nonXRE.ReplaceAllString(reString, "($1)")     // regexp group for non-x sequences
	reString = noInitialX.ReplaceAllString(reString, "()$1") // add empty prefix group if needed
	reString = noFinalX.ReplaceAllString(reString, "$1()")   // add empty suffix group if needed
	reString = strings.ReplaceAll(reString, "x", "(.*)")     // regexp group for x sequences
	//fmt.Printf("%v => /%v/\n", baseString, reString)
	re := regexp.MustCompile("^" + reString + "$")
	return re
}

func (r BaseRule) Match(input string) (MatchResult, bool) {
	if r.IsIntRule() {
		n, err := strconv.Atoi(input)
		if err != nil {
			return MatchResult{}, false
		}
		divisor := r.Divisor()
		// >> in normal rule: Divide the number by the rule's divisor and format the remainder
		right := n % divisor
		// << in normal rule: Divide the number by the rule's divisor and format the quotient
		left := n / divisor
		return MatchResult{ForwardLeft: fmt.Sprintf("%d", left), ForwardRight: fmt.Sprintf("%d", right)}, true
	}

	// String rule
	if r.stringMatchRegexp != emptyRegexp {
		r.stringMatchRegexp = buildStringMatchRegexp(r.BaseString)
	}
	m := r.stringMatchRegexp.FindStringSubmatch(input)
	if m != nil && len(m) == 4 {
		//fmt.Printf("%v => %#v\n", input, m)
		left := m[1]
		right := m[3]
		return MatchResult{ForwardLeft: left, ForwardRight: right}, true
	}
	return MatchResult{}, false
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
	ruleName := strings.TrimPrefix(ruleRef, "%")
	res, ok := g.RuleSets[ruleName]
	return res, ok
}

func NewRuleSetGroup(name string, ruleSets []RuleSet) (RuleSetGroup, error) {
	rsMap := make(map[string]RuleSet)
	for _, rs := range ruleSets {
		// sort each rule set in ascending order?
		//sort.Slice(rs.Rules, func(i, j int) bool { return rs.Rules[i].BaseInt < rs.Rules[j].BaseInt })
		rsMap[rs.Name] = rs
	}
	res := RuleSetGroup{Name: name, RuleSets: rsMap}

	for _, ruleSet := range res.RuleSets {
		for _, rule := range ruleSet.Rules {
			if rule.BaseInt != 0 && rule.BaseString != "" {
				return res, fmt.Errorf("Rule must use either BaseInt or BaseString, not both: %v", rule)
			}
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
		if r.IsIntRule() {
			n, err := strconv.Atoi(input)
			if err != nil {
				continue
			}
			if r.BaseInt <= n {
				res = r
				found = true
			} else {
				break
			}
		} else {
			if _, matches := r.Match(input); matches {
				return r, true
			}
		}
	}
	return res, found
}

func (g RuleSetGroup) Spellout(input string, ruleSetName string) (string, error) {
	if rs, ok := g.FindRuleSet(ruleSetName); ok {
		return g.spellout(input, rs)
	}
	return "", fmt.Errorf("No such rule set: %s", ruleSetName)
}

func (g RuleSetGroup) spellout(input string, ruleSet RuleSet) (string, error) {
	var err error

	matchedRule, ok := findMatchingRule(input, ruleSet)
	if !ok {
		return input, fmt.Errorf("No matching base rule for %s", input)
	}

	if fmt.Sprintf("%d", matchedRule.BaseInt) == input {
		// if n == 0 && matchedRule.BaseInt == n {
		return matchedRule.SpellOut, nil
	}

	match, ok := matchedRule.Match(input)
	if !ok {
		return input, fmt.Errorf("Couldn't get match result for rule %v, input %s", matchedRule, input)
	}

	var left, right string
	if matchedRule.RightSub == "[>>]" { // Text in brackets is omitted if the number being formatted is an even multiple of 10
		//if n%10 != 0 {
		if !strings.HasSuffix(input, "0") {
			right, err = g.spellout(match.ForwardRight, ruleSet)
			if err != nil {
				return "", err
			}
		}
	} else if matchedRule.RightSub == ">>" {
		right, err = g.spellout(match.ForwardRight, ruleSet)
		if err != nil {
			return "", err
		}
	} else if namedRuleSet, ok := g.FindRuleSet(matchedRule.RightSub); ok {
		right, err = g.spellout(match.ForwardRight, namedRuleSet)
		if err != nil {
			return "", err
		}
	} else if matchedRule.RightSub != "" {
		return "", fmt.Errorf("Unknown rule context: %s", matchedRule.RightSub)
	}

	if matchedRule.LeftSub == "<<" {
		left, err = g.spellout(match.ForwardLeft, ruleSet)
		if err != nil {
			return "", err
		}

	} else if namedRuleSet, ok := g.FindRuleSet(matchedRule.LeftSub); ok {
		left, err = g.spellout(match.ForwardLeft, namedRuleSet)
		if err != nil {
			return "", err
		}

	} else if matchedRule.LeftSub != "" {
		return "", fmt.Errorf("Unknown rule context: %s", matchedRule.LeftSub)
	}

	res := strings.TrimSpace(left + matchedRule.LeftPadding + matchedRule.SpellOut + matchedRule.RightPadding + right)
	return res, nil
}

func exp(x, y int) int {
	res := 1
	for i := 1; i <= y; i++ {
		res = res * x
	}
	return res
}

func isRuleRef(s string) bool {
	res := strings.HasPrefix(s, "%") || (s != "" && !strings.Contains(s, "<") && !strings.Contains(s, ">"))
	//fmt.Printf("%v %v\n", s, res)
	return res
}
