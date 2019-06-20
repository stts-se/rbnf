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

func TestUnmarshalXMLSV(t *testing.T) {
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

func TestReadXMLFileSV(t *testing.T) {
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

func TestRulesFromXMLFileSV(t *testing.T) {

	pack, err := RulesFromXMLFile("test_data/sv.xml")
	if err != nil {
		t.Errorf("Pain! %v", err)
	}

	if w, g := "sv", string(pack.Language); w != g {
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
	expect = "tjugo­tusen"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain for %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "200000"
	expect = "två­hundra­tusen"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain for %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "200001"
	expect = "två­hundra­tusen ett"
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

	input = "20117500"
	expect = "tjugo miljoner ett­hundra­sjutton­tusen fem­hundra"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("Stradivarius: %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "10117500"
	expect = "tio miljoner ett­hundra­sjutton­tusen fem­hundra"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("Filtsocka: %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "12345"
	expect = "12 345:e"
	res, err = pack.Spellout(input, "OrdinalRules", "digits-ordinal-masculine", true)
	if err != nil {
		t.Errorf("Filtsocka: %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "2300000007000010000"
	expect = "2 300 000 007 000 010 000:e"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-ordinal-neuter", true)
	if err != nil {
		t.Errorf("Filtsocka: %s! %v", input, err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

}

func TestRulesFromXMLFileDE(t *testing.T) {

	pack, err := RulesFromXMLFile("test_data/de.xml")
	if err != nil {
		t.Errorf("Pain! %v", err)
	}

	if w, g := "de", string(pack.Language); w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	ruleSetGroups := pack.RuleSetGroups
	if len(ruleSetGroups) == 0 {
		t.Errorf("Noooo!")
	}

	if w, g := "SpelloutRules", ruleSetGroups[0].Name; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	var input, expect, res string

	input = "21"
	expect = "ein­und­zwanzig"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "48"
	expect = "acht­und­vierzig"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "2748"
	expect = "zwei­tausend­sieben­hundert­acht­und­vierzig"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	//
	input = "13000"
	expect = "dreizehn­tausend"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

}

func TestRulesFromXMLFileFR(t *testing.T) {

	pack, err := RulesFromXMLFile("test_data/fr.xml")
	if err != nil {
		t.Errorf("Pain! %v", err)
	}

	if w, g := "fr", string(pack.Language); w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	ruleSetGroups := pack.RuleSetGroups
	if len(ruleSetGroups) == 0 {
		t.Errorf("Noooo!")
	}

	if w, g := "SpelloutRules", ruleSetGroups[0].Name; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	var input, expect, res string

	input = "78"
	expect = "soixante-dix-huit"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "8765"
	expect = "huit mille sept cent soixante-cinq"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "485"
	expect = "quatre cent quatre-vingt-cinq"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "435"
	expect = "quatre cent trente-cinq"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

}

func TestRulesFromXMLFileTA(t *testing.T) {

	pack, err := RulesFromXMLFile("test_data/ta.xml")
	if err != nil {
		t.Errorf("Pain! %v", err)
	}

	if w, g := "ta", string(pack.Language); w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	ruleSetGroups := pack.RuleSetGroups
	if len(ruleSetGroups) == 0 {
		t.Errorf("Noooo!")
		return
	}

	if w, g := "SpelloutRules", ruleSetGroups[0].Name; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	var input, expect, res string

	input = "78"
	expect = "எழுபது எட்டு"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "8765"
	expect = "எட்டு ஆயிரம் எழுநூறு அறுபது ஐந்து"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "485"
	expect = "நாநூறூ எண்பது ஐந்து"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

	input = "935"
	expect = "தொள்ளாயிரம் முப்பது ஐந்து"
	res, err = pack.Spellout(input, "SpelloutRules", "spellout-numbering", false)
	if err != nil {
		t.Errorf("P-P-Pure Pain! %v", err)
	} else if res != expect {
		t.Errorf("wanted %s, got %s", expect, res)
	}

}
