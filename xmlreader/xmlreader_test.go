package xmlreader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
)

// just to have not to de-import fmt
func skunk() {
	fmt.Println()
}

func TestStruct(t *testing.T) {
	ldml := Ldml{}
	_ = ldml.Identity
}

func TestUnmarshalXML(t *testing.T) {
	// depends on file sv.xml in current dir

	bytes, err := ioutil.ReadFile("sv.xml")
	if err != nil {
		t.Errorf("Sorry! : %v", err)
	}

	ldml := Ldml{}
	err = xml.Unmarshal(bytes, &ldml)
	if err != nil {
		t.Errorf("Sad! : %v", err)
	}

	//fmt.Printf("%#v\n", ldml.Rbnf.RulesetGrouping[0])
}

func TestReadXMLFile(t *testing.T) {
	// depends on file sv.xml in current dir

	ldml, err := readXMLFile("sv.xml")
	if err != nil {
		t.Errorf("Sob! : %v", err)
	}

	//fmt.Printf("%#v\n")

	if w, g := "sv", ldml.Identity.Language.Attrtype; w != g {
		t.Errorf("wanted %s got %s", w, g)
	}

}

func TestRulesFromXMLFile(t *testing.T) {

	lang, ruleSetGroups, err := RulesFromXMLFile("sv.xml")
	if err != nil {
		t.Errorf("Pain! %v", err)
	}

	if w, g := "sv", lang; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	if len(ruleSetGroups) == 0 {
		t.Errorf("Noooo!")
	}

	if w, g := "SpelloutRules", ruleSetGroups[0].Name; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	sor := ruleSetGroups[0]

	fmt.Printf("%#v\n", sor)

	if w, g := "spellout-numbering-t", sor.RuleSets["spellout-numbering-t"].Name; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}
}
