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
	// depends on file sv.xml in test_data dir

	bytes, err := ioutil.ReadFile("test_data/sv.xml")
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
	// depends on file sv.xml in test_data dir

	ldml, err := readXMLFile("test_data/sv.xml")
	if err != nil {
		t.Errorf("Sob! : %v", err)
	}

	//fmt.Printf("%#v\n")

	if w, g := "sv", ldml.Identity.Language.Attrtype; w != g {
		t.Errorf("wanted %s got %s", w, g)
	}

}

func TestRulesFromXMLFile(t *testing.T) {

	pack, err := RulesFromXMLFile("test_data/sv.xml")
	if err != nil {
		t.Errorf("Pain! %v", err)
	}

	if w, g := "sv", pack.Language; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	ruleSetGroups := pack.RuleSetGroups
	if len(ruleSetGroups) == 0 {
		t.Errorf("Noooo!")
	}

	if w, g := "SpelloutRules", ruleSetGroups[0].Name; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	sor := ruleSetGroups[0]

	//fmt.Printf("%#v\n", sor)

	if w, g := "spellout-numbering-t", sor.RuleSets["spellout-numbering-t"].Name; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	var input, expect, res string

	//
	input = "10"
	expect = "tio"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "20"
	expect = "tjugo"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain for %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "20000"
	expect = "tjugo tusen"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain for %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "200000"
	expect = "två hundra tusen"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain for %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "200001"
	expect = "två hundra tusen ett"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain for %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "2000000"
	expect = "två miljoner"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain for %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

}
