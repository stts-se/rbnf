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

type Base struct {

	// Int base
	Int   int
	Radix int // only used for Int base

	// String base
	String string
}

func (b Base) ToString() string {
	if b.IsInt() {
		return fmt.Sprintf("%d (%d)", b.Int, b.Radix)
	}
	return b.String
}

type BaseRule struct {
	Base         Base
	LeftSub      string
	LeftPadding  string
	SpellOut     []string
	RightPadding string
	RightSub     string

	stringMatchRegexp *regexp.Regexp
}

func NewIntRule(baseInt int, leftSub, leftPadding string, spellOut string, rightPadding, rightSub string, radix int) BaseRule {
	return BaseRule{
		Base:         Base{Int: baseInt, Radix: radix},
		LeftSub:      leftSub,
		LeftPadding:  leftPadding,
		SpellOut:     []string{spellOut},
		RightPadding: rightPadding,
		RightSub:     rightSub,
	}
}
func NewStringRule(baseString string, leftSub, leftPadding string, spellOut string, rightPadding, rightSub string) BaseRule {
	return BaseRule{
		Base:              Base{String: baseString},
		LeftSub:           leftSub,
		LeftPadding:       leftPadding,
		SpellOut:          []string{spellOut},
		RightPadding:      rightPadding,
		RightSub:          rightSub,
		stringMatchRegexp: buildStringMatchRegexp(baseString),
	}
}

func (r BaseRule) String() string {
	if r.Base.IsInt() {
		return fmt.Sprintf("%s => '%s%s%s%s%s'", r.Base.ToString(), r.LeftSub, r.LeftPadding, strings.Join(r.SpellOut, ""), r.RightPadding, r.RightSub)
	}
	return fmt.Sprintf("%s => '%s%s%s%s%s'", r.Base.ToString(), r.LeftSub, r.LeftPadding, strings.Join(r.SpellOut, ""), r.RightPadding, r.RightSub)
}

func (b Base) IsInt() bool {
	return b.String == ""
}

func (b Base) Divisor() int {

	if !b.IsInt() {
		panic("invalid call to Base.Divisor for String type Base")
	}

	/** http://icu-project.org/apiref/icu4c/classRuleBasedNumberFormat.html
	To calculate the divisor, let [...] the exponent be the highest exponent of the radix that yields a result less than or equal to the base value.
	If the exponent is positive or 0, the divisor is the radix raised to the power of the exponent; otherwise, the divisor is 1.
	*/

	//for rad >= 0
	//exponent : the highest exponent of the radix that is less than or equal to the base value
	//divisor: radix^exponent
	var exponent, divisor int
	for i := 1; exp(b.Radix, i) <= b.Int; i++ {
		exponent = i
	}
	if exponent >= 0 {
		divisor = exp(b.Radix, exponent)
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
	if r.Base.IsInt() {
		n, err := strconv.Atoi(input)
		if err != nil {
			return MatchResult{}, false
		}
		divisor := r.Base.Divisor()
		// >> in normal rule: Divide the number by the rule's divisor and format the remainder
		right := n % divisor
		// << in normal rule: Divide the number by the rule's divisor and format the quotient
		left := n / divisor
		return MatchResult{ForwardLeft: fmt.Sprintf("%d", left), ForwardRight: fmt.Sprintf("%d", right)}, true
	}

	// String rule
	if r.stringMatchRegexp != emptyRegexp || r.stringMatchRegexp == nil {
		r.stringMatchRegexp = buildStringMatchRegexp(r.Base.String)
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
	Language string
	//RuleSetGroups map[string]RuleSetGroup
	RuleSetGroups []RuleSetGroup
}

func (p RulePackage) Spellout(input string, groupName string, ruleSetName string) (string, error) {
	// if g, ok := p.RuleSetGroups[groupName]; ok {
	// 	res, err := g.Spellout(input, ruleSetName)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return res, nil
	// }
	for _, g := range p.RuleSetGroups {
		if g.Name == groupName {
			res, err := g.Spellout(input, ruleSetName)
			if err != nil {
				return "", err
			}
			return res, nil
		}
	}
	return "", fmt.Errorf("No such rule set group: %s", groupName)
}

type RuleSetGroup struct {
	Name     string
	RuleSets map[string]RuleSet
}

func (g RuleSetGroup) FindRuleSet(ruleRef string) (RuleSet, bool) {
	ruleName := ruleRef
	ruleName = strings.TrimPrefix(ruleName, "<")
	ruleName = strings.TrimPrefix(ruleName, "%")
	ruleName = strings.TrimPrefix(ruleName, "%")
	ruleName = strings.TrimSuffix(ruleName, "<")
	res, ok := g.RuleSets[ruleName]
	return res, ok
}

func (g RuleSetGroup) FindSpelloutRuleSet(ruleRef string) (RuleSet, bool) {
	ruleName := ruleRef
	ruleName = strings.TrimPrefix(ruleName, "=%")
	ruleName = strings.TrimSuffix(ruleName, "=")
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
			if rule.Base.Int != 0 && rule.Base.String != "" {
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
			for _, sp := range rule.SpellOut {
				if isSpelloutRuleRef(sp) {
					if _, ok := res.FindSpelloutRuleSet(sp); !ok {
						return res, fmt.Errorf("No such rule set: %s", sp)
					}
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
		if r.Base.IsInt() {
			n, err := strconv.Atoi(input)
			if err != nil {
				continue
			}
			if r.Base.Int <= n {
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

func (g RuleSetGroup) expandSpellouts(r BaseRule) (string, error) {
	res := []string{}
	for _, sp := range r.SpellOut {
		if isSpelloutRuleRef(sp) {
			rs, ok := g.FindSpelloutRuleSet(sp)
			if !ok {
				return "", fmt.Errorf("No such rule set: %s", sp)
			}
			spelled, err := g.spellout(r.Base.ToString(), rs)
			if err != nil {
				return "", nil
			}
			res = append(res, spelled)
		} else {
			res = append(res, sp)
		}
	}
	// TODO: Call rule defs too
	return strings.Join(res, ""), nil
}

func (g RuleSetGroup) spellout(input string, ruleSet RuleSet) (string, error) {
	var err error

	matchedRule, ok := findMatchingRule(input, ruleSet)
	if !ok {
		return input, fmt.Errorf("No matching base rule for %s", input)
	}

	if fmt.Sprintf("%d", matchedRule.Base.Int) == input {
		//if n, err := strconv.Atoi(input); err != nil && n == 0 && matchedRule.Base.Int == n {
		sp, err := g.expandSpellouts(matchedRule)
		if err != nil {
			return "", err
		}
		return sp, nil
	}

	match, ok := matchedRule.Match(input)
	if !ok {
		return input, fmt.Errorf("Couldn't get match result for rule %v, input %s", matchedRule, input)
	}

	var left, right string
	if matchedRule.RightSub == "[>>]" { // ??? Text in brackets is omitted if the number being formatted is an even multiple of 10
		n, err := strconv.Atoi(input)
		//fmt.Println(n)
		omit := n%10 == 0
		omit = false
		if err != nil || !omit {
			//if matchedRule.Base.IsInt() && matchedRule.Base.Int%10 != 0 {
			//if n%10 != 0 {
			///if !strings.HasSuffix(input, "0") {
			right, err = g.spellout(match.ForwardRight, ruleSet)
			if err != nil {
				return "", err
			}
			//}
		}
	} else if matchedRule.RightSub == ">>" {
		//fmt.Printf("xx %#v\t%v\n", matchedRule, match.ForwardRight)
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
		return "", fmt.Errorf("Unknown rule context: '%s'", matchedRule.RightSub)
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
		return "", fmt.Errorf("Unknown rule context: '%s'", matchedRule.LeftSub)
	}

	spell, err := g.expandSpellouts(matchedRule)
	if err != nil {
		return "", err
	}
	res := strings.TrimSpace(left + matchedRule.LeftPadding + spell + matchedRule.RightPadding + right)
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

func isSpelloutRuleRef(s string) bool {
	res := strings.HasPrefix(s, "=%") && strings.HasSuffix(s, "=")
	//fmt.Printf("%v %v\n", s, res)
	return res
}
