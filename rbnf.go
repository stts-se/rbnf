package rbnf

import (
	"fmt"
	"os"
	"regexp"
	// "sort"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type RuleSet struct {
	Name    string
	Rules   []BaseRule
	Private bool
}

type Base struct {

	// Int base
	Int   int
	Radix int // only used for Int base

	// String base
	String            string
	StringMatchRegexp *Regexp
}

func NewBaseString(s string) Base {
	return Base{
		String:            s,
		StringMatchRegexp: buildStringMatchRegexp(s),
	}

}

func NewBaseInt(n int, radix int) Base {
	return Base{
		Int:   n,
		Radix: radix,
	}

}

type Language string

func (b Base) ToString() string {
	if b.IsInt() {
		return fmt.Sprintf("%d (%d)", b.Int, b.Radix)
	}
	return b.String
}

func (b Base) Value() string {
	if b.IsInt() {
		return fmt.Sprintf("%d", b.Int)
	}
	return b.String
}

type Formatter struct {
	printer     *message.Printer
	format      string
	initialized bool
}

func (f Formatter) String() string {
	return f.format
}

type Sub struct {
	Optional  bool
	Orth      string
	RuleRef   string
	Formatter Formatter
	Operation string
}

func ParseSub(sub string, lang Language) Sub {
	input := sub
	res := Sub{}
	if len([]rune(sub)) == 0 {
		return res
	}
	if strings.HasPrefix(sub, "[") && strings.HasSuffix(sub, "]") {
		res.Optional = true
		sub = strings.TrimPrefix(strings.TrimSuffix(sub, "]"), "[")
	}
	if len([]rune(sub)) == 0 {
		return res
	}
	firstChar := string([]rune(sub)[0])
	if (firstChar == ">" || firstChar == "<" || firstChar == "=") && strings.HasSuffix(sub, firstChar) {
		res.Operation = firstChar + firstChar
		ref := strings.TrimPrefix(strings.TrimSuffix(sub, firstChar), firstChar)
		if strings.HasPrefix(ref, "#") || (res.Operation != "" && strings.Contains(ref, "0")) {
			p := message.NewPrinter(language.Make(string(lang)))
			res.Formatter = Formatter{printer: p, format: ref, initialized: true}
		} else {
			res.RuleRef = ref
		}
	} else {
		res.Orth = sub
	}
	if res.String() != input {
		msg := fmt.Sprintf("wtf! %s != %s (%#v)", input, res.String(), res)
		panic(msg)
	}
	return res
}

func (sub Sub) String() string {
	res := ""
	if sub.Orth != "" {
		res = sub.Orth
	} else if sub.RuleRef != "" {
		res = sub.RuleRef
	} else if sub.Formatter.initialized {
		res = sub.Formatter.String()
	}
	if sub.Operation != "" {
		op := []rune(sub.Operation)[0]
		res = fmt.Sprintf("%s%s%s", string(op), res, string(op))
	}
	if sub.Optional {
		res = "[" + res + "]"
	}
	return res
}

func (sub Sub) IsRuleRef() bool {
	return sub.RuleRef != "" && strings.HasPrefix(sub.RuleRef, "%")
}

func (sub Sub) IsFormatRef() bool {
	return sub.Formatter.initialized
}

func (sub Sub) IsError() bool {
	return sub.Orth == "ERROR"
}

func (sub Sub) Validate() error {
	if sub.Orth != "" && sub.RuleRef != "" {
		return fmt.Errorf("Orth and RuleRef cannot both be instantiated")
	}
	if sub.Orth != "" && sub.Operation != "" {
		return fmt.Errorf("Orth and Operation cannot both be instantiated")
	}
	return nil
}

type Regexp struct {
	initialised bool
	re          *regexp.Regexp
}

type BaseRule struct {
	Base Base
	Subs []Sub
}

func NewIntRule(lang Language, baseInt int, radix int, subs ...string) BaseRule {
	subSubs := []Sub{}
	for _, s := range subs {
		sub := ParseSub(s, lang)
		subSubs = append(subSubs, sub)
	}
	return BaseRule{
		Base: Base{Int: baseInt, Radix: radix},
		Subs: subSubs,
	}
}
func NewStringRule(lang Language, baseString string, subs ...string) BaseRule {
	subSubs := []Sub{}
	for _, s := range subs {
		sub := ParseSub(s, lang)
		subSubs = append(subSubs, sub)
	}
	return BaseRule{
		Base: Base{String: baseString, StringMatchRegexp: buildStringMatchRegexp(baseString)},
		Subs: subSubs,
	}
}

func (r *BaseRule) String() string {
	subs := []string{}
	for _, sub := range r.Subs {
		subs = append(subs, sub.String())
	}
	return fmt.Sprintf("%v => '%s'", r.Base.ToString(), strings.Join(subs, ""))
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

// TODO: this is sooo ugly -- can it be done better?
func buildStringMatchRegexp(baseString string) *Regexp {
	reString := baseString
	reString = regexpEscape(reString)                        // escape special chars in the BaseString
	reString = nonXRE.ReplaceAllString(reString, "($1)")     // regexp group for non-x sequences
	reString = noInitialX.ReplaceAllString(reString, "()$1") // add empty prefix group if needed
	reString = noFinalX.ReplaceAllString(reString, "$1()")   // add empty suffix group if needed
	reString = strings.ReplaceAll(reString, "x", "(.*)")     // regexp group for x sequences
	//fmt.Printf("%v => /%v/\n", baseString, reString)
	re := regexp.MustCompile("^" + reString + "$")
	return &Regexp{initialised: true, re: re}
}

func (r *BaseRule) Match(input string) (MatchResult, bool) {
	// A) Int rule
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

	// B) String rule

	// TODO
	var fasterMatch = true

	if fasterMatch {
		switch r.Base.String {
		case "x.x":
			pts := strings.Split(input, ".")
			if len(pts) == 2 {
				res := MatchResult{ForwardLeft: pts[0], ForwardRight: pts[1]}
				return res, true
			}
		case "x,x":
			pts := strings.Split(input, ",")
			if len(pts) == 2 {
				res := MatchResult{ForwardLeft: pts[0], ForwardRight: pts[1]}
				return res, true
			}
		case "-x":
			pts := strings.Split(input, "-")
			if len(pts) == 2 {
				res := MatchResult{ForwardLeft: pts[0], ForwardRight: pts[1]}
				return res, true
			}
		case "x%":
			pts := strings.Split(input, "%")
			if len(pts) == 2 {
				res := MatchResult{ForwardLeft: pts[0], ForwardRight: pts[1]}
				return res, true
			}
			//case "#,##": return MatchResult{}, false
		default:
			return MatchResult{}, false
		}
	} else {
		if !r.Base.StringMatchRegexp.initialised {
			r.Base.StringMatchRegexp = buildStringMatchRegexp(r.Base.String)
		}
		//fmt.Println("RULE AND REGEXP:", r, r.Base.StringMatchRegexp.re)
		m := r.Base.StringMatchRegexp.re.FindStringSubmatch(input)
		if m != nil && len(m) == 4 {
			//fmt.Printf("%v => %#v\n", input, m)
			left := m[1]
			right := m[3]
			return MatchResult{ForwardLeft: left, ForwardRight: right}, true
		}
	}
	return MatchResult{}, false
}

type MatchResult struct {
	ForwardLeft  string
	ForwardRight string
}

type RulePackage struct {
	Language      Language
	RuleSetGroups []RuleSetGroup
	Debug         bool
}

func (r *RulePackage) Spellout(input string, groupName string, ruleSetName string, debug bool) (string, error) {
	for _, g := range r.RuleSetGroups {
		if g.Name == groupName {
			res, err := g.Spellout(input, ruleSetName, debug)
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
	Language Language
	RuleSets map[string]RuleSet
}

func (g RuleSetGroup) FindRuleSet(ruleRef string) (RuleSet, bool) {
	ruleName := ruleRef
	ruleName = strings.TrimPrefix(ruleName, "%")
	ruleName = strings.TrimPrefix(ruleName, "%")
	res, ok := g.RuleSets[ruleName]
	return res, ok
}

func NewRulePackage(lang Language, ruleSetGroups []RuleSetGroup, debug bool) (RulePackage, error) {
	res := RulePackage{Language: lang, Debug: debug, RuleSetGroups: ruleSetGroups}
	for _, g := range res.RuleSetGroups {
		if g.Language != res.Language {
			return res, fmt.Errorf("Language for rule set group %s does not match package language: %s / %s", g.Name, res.Language, g.Language)
		}
		for _, ruleSet := range g.RuleSets {
			for _, rule := range ruleSet.Rules {
				if rule.Base.Int != 0 && rule.Base.String != "" {
					return res, fmt.Errorf("Rule must use either BaseInt or BaseString, not both: %v", rule)
				}
				for _, sub := range rule.Subs {
					if sub.IsRuleRef() {
						if _, ok := g.FindRuleSet(sub.RuleRef); !ok {
							return res, fmt.Errorf("No such rule set: %s", sub)
						}
					}
				}
			}
		}
	}
	return res, nil
}

func (g RuleSetGroup) Validate() error {
	for _, ruleSet := range g.RuleSets {
		for _, rule := range ruleSet.Rules {
			if rule.Base.Int != 0 && rule.Base.String != "" {
				return fmt.Errorf("Rule must use either BaseInt or BaseString, not both: %v", rule)
			}
			for _, sub := range rule.Subs {
				if sub.IsRuleRef() {
					if _, ok := g.FindRuleSet(sub.RuleRef); !ok {
						return fmt.Errorf("No such rule set: %s", sub)
					}
				}
			}
		}
	}
	return nil
}

func NewRuleSetGroup(name string, lang Language, ruleSets []RuleSet) (RuleSetGroup, error) {
	rsMap := make(map[string]RuleSet)
	for _, rs := range ruleSets {
		// sort each rule set in ascending order?
		//sort.Slice(rs.Rules, func(i, j int) bool { return rs.Rules[i].BaseInt < rs.Rules[j].BaseInt })
		rsMap[rs.Name] = rs
	}
	res := RuleSetGroup{Name: name, Language: lang, RuleSets: rsMap}

	err := res.Validate()

	return res, err
}

func findMatchingRule(input string, ruleSet RuleSet) (BaseRule, bool) {
	var res BaseRule
	var found = false

	n, err := strconv.Atoi(input)

	for _, r := range ruleSet.Rules {
		if r.Base.IsInt() {
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

// details here: http://www.icu-project.org/applets/icu4j/4.1/docs-4_1_1/com/ibm/icu/text/DecimalFormat.html
func format(input string, formatter Formatter, debug bool) (string, error) {
	if debug {
		fmt.Fprintf(os.Stderr, "[rbnf.format] Input:%s Fmt:%s\n", input, formatter.format)
	}
	//var p = message.NewPrinter(language.Make(lang))
	var numeric interface{}
	var err error
	numeric, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		numeric, err = strconv.ParseFloat(input, 64)
		if err != nil {
			return input, err
		}
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[rbnf.format] numeric: %v\n", numeric)
	}
	res := formatter.printer.Sprint(numeric)
	return res, nil
}

func (g *RuleSetGroup) Spellout(input string, ruleSetName string, debug bool) (string, error) {
	if rs, ok := g.FindRuleSet(ruleSetName); ok {
		res, err := g.spellout(input, rs, debug)
		return strings.TrimSpace(res), err

	}
	return "", fmt.Errorf("No such rule set: %s", ruleSetName)
}

func (g *RuleSetGroup) spellout(input string, ruleSet RuleSet, debug bool) (string, error) {
	matchedRule, ok := findMatchingRule(input, ruleSet)
	if !ok {
		err := fmt.Errorf("No matching base rule for %s", input)
		if debug {
			fmt.Fprintf(os.Stderr, "[rbnf] %v : rule set: %#v\n", err, ruleSet)
		}
		return input, err
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[rbnf] Input %v\n", input)
		fmt.Fprintf(os.Stderr, "[rbnf] Matched rule %#v (from rule set %s)\n", matchedRule, ruleSet.Name)
	}

	match, ok := matchedRule.Match(input)
	if !ok {
		return input, fmt.Errorf("Couldn't get match result for rule %v, input %s", matchedRule, input)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[rbnf] Match result: %#v\n", match)
	}

	var subs = []string{}
	for _, sub := range matchedRule.Subs {
		if debug {
			fmt.Fprintf(os.Stderr, "[rbnf] Current sub: %#v\n", sub)
		}
		// http://www.icu-project.org/applets/icu4j/4.1/docs-4_1_1/com/ibm/icu/text/RuleBasedNumberFormat.html
		// Omit the optional text if the number is an even multiple of the rule's divisor
		if sub.Optional {
			if inputInt, err := strconv.Atoi(input); err == nil && matchedRule.Base.IsInt() {
				if debug {
					fmt.Fprintf(os.Stderr, "[rbnf.optional] matchedRule divisor %v\n", matchedRule.Base.Divisor())
				}
				omit := inputInt%matchedRule.Base.Divisor() == 0
				if debug {
					fmt.Fprintf(os.Stderr, "[rbnf.optional] Omit %v\n", omit)
				}
				if omit {
					continue
				}
			}
		}

		if debug {
			fmt.Fprintf(os.Stderr, "[rbnf] Accumulated subs: %#v\n", subs)
		}
		if sub.IsFormatRef() {
			if sub.Operation == ">>" {
				spelled, err := format(match.ForwardRight, sub.Formatter, debug)
				if err != nil {
					return "", err
				}
				subs = append(subs, spelled)
			} else if sub.Operation == "<<" {
				spelled, err := format(match.ForwardLeft, sub.Formatter, debug)
				if err != nil {
					return "", err
				}
				subs = append(subs, spelled)
			} else if sub.Operation == "==" {
				spelled, err := format(input, sub.Formatter, debug)
				if err != nil {
					return "", err
				}
				subs = append(subs, spelled)
			} else {
				return input, fmt.Errorf("unknown operation for sub %#v : %s", sub, sub.Operation)
			}
		} else if namedRuleSet, ok := g.FindRuleSet(sub.RuleRef); ok {
			if sub.Operation == ">>" {
				spelled, err := g.spellout(match.ForwardRight, namedRuleSet, debug)
				if err != nil {
					return "", err
				}
				subs = append(subs, spelled)
			} else if sub.Operation == "<<" {
				spelled, err := g.spellout(match.ForwardLeft, namedRuleSet, debug)
				if err != nil {
					return "", err
				}
				subs = append(subs, spelled)
			} else if sub.Operation == "==" {
				spelled, err := g.spellout(input, namedRuleSet, debug)
				if err != nil {
					return "", err
				}
				subs = append(subs, spelled)
			} else {
				return input, fmt.Errorf("unknown operation for sub %#v : %s", sub, sub.Operation)
			}
		} else if sub.IsRuleRef() {
			return input, fmt.Errorf("unknown rule set ref %s in rule sub %s", sub.RuleRef, sub)
		} else if sub.IsError() {
			return input, fmt.Errorf("Internal rule error for input string %s (pre-defined in rule)", input)
		} else if sub.Operation == ">>" {
			spelled, err := g.spellout(match.ForwardRight, ruleSet, debug)
			if err != nil {
				return "", err
			}
			subs = append(subs, spelled)
		} else if sub.Operation == "<<" {
			spelled, err := g.spellout(match.ForwardLeft, ruleSet, debug)
			if err != nil {
				return "", err
			}
			subs = append(subs, spelled)
		} else if sub.Orth != "" {
			subs = append(subs, sub.Orth)
		}
	}

	for i, sub := range subs {
		subs[i] = strings.Trim(sub, "'")
	}
	res := strings.Join(subs, "")
	//res = strings.TrimSpace(res)       // trim space  -- ga 120.000 doesn't work with trimspace here
	res = strings.Replace(res, "  ", " ", -1)
	if res == "" {
		if debug {
			fmt.Fprintf(os.Stderr, "empty output for input string %s", input)
		}
		if !ruleSet.Private {
			return input, fmt.Errorf("empty output for input string %s", input)
		}
	}
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
	res := strings.HasPrefix(s, "%") // || (s != "" && !strings.Contains(s, "<") && !strings.Contains(s, ">"))
	return res
}

func isSpelloutRuleRef(s string) bool {
	res := strings.HasPrefix(s, "=%") && strings.HasSuffix(s, "=")
	//fmt.Printf("%v %v\n", s, res)
	return res
}
