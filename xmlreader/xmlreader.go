/** Package xmlreader contains a parser for the CLDR RBNF format https://github.com/unicode-org/cldr/tree/master/common/rbnf

License for CLDR: https://github.com/unicode-org/cldr/blob/master/ICU-LICENSE

*/
package xmlreader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/stts-se/rbnf"
	"github.com/stts-se/rbnf/lexer"
)

var Verb = false

func readXMLFile(fn string) (Ldml, error) {
	res := Ldml{}

	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return res, fmt.Errorf("failed to read XML file : %v", err)
	}

	err = xml.Unmarshal(bytes, &res)
	if err != nil {
		return res, fmt.Errorf("failed to process XML file : %v", err)
	}

	return res, nil
}

func readXMLURL(url string) (Ldml, error) {
	res := Ldml{}

	resp, err := http.Get(url)
	if err != nil {
		return res, fmt.Errorf("failed to read URL : %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("failed to read URL : %v", resp.Status)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, fmt.Errorf("failed to read XML file : %v", err)
	}

	err = xml.Unmarshal(bytes, &res)
	if err != nil {
		return res, fmt.Errorf("failed to process XML file : %v", err)
	}

	return res, nil
}

func replaceChars(s string) string {
	s = strings.Replace(s, "→", ">", -1)
	s = strings.Replace(s, "←", "<", -1)
	s = strings.Replace(s, "−", "-", -1)
	//s = strings.Replace(s, "\u00ad", "", -1) // soft hyphen
	return s
}

var threeArrows = regexp.MustCompile("(→%+[a-z-]*→[a-z-]*→|←%+[a-z-]*←[a-z-]*←)")

func unsupportedRuleFormat(rFmt string) bool {
	return strings.Contains(rFmt, "ignorable") ||
		//strings.Contains(rFmt, "$") ||
		strings.Contains(rFmt, "→→→") ||
		threeArrows.MatchString(rFmt)
}

func convertRuleSet(rs *Ruleset, lang string) (rbnf.RuleSet, error) {
	var res rbnf.RuleSet
	res.Name = rs.Attrtype
	if rs.Attraccess == "private" {
		res.Private = true
	}
	for _, r := range rs.Rbnfrule {
		//fmt.Printf("RULE %#v\n", r)
		rule := rbnf.BaseRule{}
		baseNum, err := strconv.Atoi(strings.Replace(r.Attrvalue, ",", "", -1))
		if err == nil { // numeric rule
			// TODO test
			radix := 10 // Default radix
			if r.Attrradix != "" {
				radix, err = strconv.Atoi(strings.Replace(r.Attrradix, ",", "", -1))
				if err != nil {
					return res, fmt.Errorf("failed to convert radix : %v\n", err)
				}

			}
			rule.Base = rbnf.NewBaseInt(baseNum, radix)
		} else { // non-numeric rule
			rule.Base = rbnf.NewBaseString(r.Attrvalue)
		}

		if unsupportedRuleFormat(r.String) {
			if Verb {
				log.Printf("[xmlreader] skipping unsupported rule format: %#v", r)
				// } else {
				// 	log.Printf("[xmlreader] skipping unsupported rule format: %#v", r.String)
			}
			continue
		}

		lex := lexer.Lex(r.String)
		err = lex.Run()

		if err != nil {
			err = fmt.Errorf("parse failed for '%v' : %v", r, err)
			if Verb {
				log.Printf("[xmlreader] %v", err)
			}
			return res, err

		}
		for _, i := range lex.Result() {
			sub, err := rbnf.ParseSub(replaceChars(i), rbnf.Language(lang))
			if err != nil {
				if Verb {
					log.Printf("[xmlreader] %v", err)
				}
				return res, err
			}
			rule.Subs = append(rule.Subs, sub)
		}
		// if Verb {
		// 	log.Printf("PARSED RULE\t%#v\t%#v\t%#v\t%#v\n", res.Name, r.String, rule, rule)
		// }
		res.Rules = append(res.Rules, rule)
	}

	return res, nil
}

func convertGroup(g *RulesetGrouping, lang string) (string, []rbnf.RuleSet, error) {
	var res []rbnf.RuleSet
	name := g.Attrtype
	if strings.TrimSpace(name) == "" {
		return "", res, fmt.Errorf("rule set grouping lacks type attribute value")
	}

	for _, rs := range g.Ruleset {
		rbnfRuleSet, err := convertRuleSet(rs, lang)
		if err != nil {
			return name, res, fmt.Errorf("failed to convert rule set : %v", err)
			//fmt.Fprintf(os.Stderr, "skipping rule set '%s' : %v\n", rbntRuleSet.Name, err)
			//continue
		}
		if len(rbnfRuleSet.Rules) > 0 {
			res = append(res, rbnfRuleSet)
		}
	}

	if len(res) > 0 {
		return name, res, nil
	}
	return name, res, fmt.Errorf("no rule sets for rule set group %s", name)
}

func rulesFromLdml(ldml Ldml, lang string) ([]rbnf.RuleSetGroup, error) {
	//var res rbnf.RuleSetGroup
	res := []rbnf.RuleSetGroup{}
	// name of whole file

	groups := ldml.Rbnf.RulesetGrouping
	if len(groups) == 0 {
		return res, fmt.Errorf("empty RulsetGrouping")
	}

	var rbnfGroups []rbnf.RuleSetGroup
	for _, g := range groups {
		name, ruleSet, err := convertGroup(g, lang)
		if err != nil {
			return res, fmt.Errorf("failed to convert rule group : %v", err)
			//fmt.Fprintf(os.Stderr, "skipping rule group '%s' : %v", name, err)
			//continue
		}
		group, err := rbnf.NewRuleSetGroup(name, rbnf.Language(lang), ruleSet)
		if err != nil {

			//fmt.Printf("%#v\n", group)
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

func RulesFromXMLFile(fn string) (rbnf.RulePackage, error) {
	var lang string

	ldml, err := readXMLFile(fn)
	if err != nil {
		return rbnf.RulePackage{}, fmt.Errorf("RulesFromXMLFile: %v", err)
	}
	lang = ldml.Identity.Language.Attrtype

	groups, err := rulesFromLdml(ldml, lang)
	if err != nil {
		return rbnf.RulePackage{}, fmt.Errorf("RulesFromXMLFile: %v", err)
	}

	return rbnf.NewRulePackage(rbnf.Language(lang), groups, false)
}

func RulesFromXMLURL(url string) (rbnf.RulePackage, error) {
	var lang string

	ldml, err := readXMLURL(url)
	if err != nil {
		return rbnf.RulePackage{}, fmt.Errorf("RulesFromXMLURL: %v", err)
	}
	lang = ldml.Identity.Language.Attrtype

	groups, err := rulesFromLdml(ldml, lang)
	if err != nil {
		return rbnf.RulePackage{}, fmt.Errorf("RulesFromXMLURL: %v", err)
	}

	return rbnf.NewRulePackage(rbnf.Language(lang), groups, false)
}
