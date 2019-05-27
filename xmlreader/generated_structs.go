type Identity struct {
	XMLName xml.Name `xml:"identity,omitempty" json:"identity,omitempty"`
	Language *Language `xml:"language,omitempty" json:"language,omitempty"`
	Version *Version `xml:"version,omitempty" json:"version,omitempty"`
}

type Language struct {
	XMLName xml.Name `xml:"language,omitempty" json:"language,omitempty"`
	Attrtype string`xml:"type,attr"  json:",omitempty"`
}

type Ldml struct {
	XMLName xml.Name `xml:"ldml,omitempty" json:"ldml,omitempty"`
	Identity *Identity `xml:"identity,omitempty" json:"identity,omitempty"`
	Rbnf *Rbnf `xml:"rbnf,omitempty" json:"rbnf,omitempty"`
}

type Rbnf struct {
	XMLName xml.Name `xml:"rbnf,omitempty" json:"rbnf,omitempty"`
	RulesetGrouping []*RulesetGrouping `xml:"rulesetGrouping,omitempty" json:"rulesetGrouping,omitempty"`
}

type Rbnfrule struct {
	XMLName xml.Name `xml:"rbnfrule,omitempty" json:"rbnfrule,omitempty"`
	Attrradix string`xml:"radix,attr"  json:",omitempty"`
	Attrvalue string`xml:"value,attr"  json:",omitempty"`
	string string `xml:",chardata" json:",omitempty"`
}

type Ruleset struct {
	XMLName xml.Name `xml:"ruleset,omitempty" json:"ruleset,omitempty"`
	Attraccess string`xml:"access,attr"  json:",omitempty"`
	Attrtype string`xml:"type,attr"  json:",omitempty"`
	Rbnfrule []*Rbnfrule `xml:"rbnfrule,omitempty" json:"rbnfrule,omitempty"`
}

type RulesetGrouping struct {
	XMLName xml.Name `xml:"rulesetGrouping,omitempty" json:"rulesetGrouping,omitempty"`
	Attrtype string`xml:"type,attr"  json:",omitempty"`
	Ruleset []*Ruleset `xml:"ruleset,omitempty" json:"ruleset,omitempty"`
}

type Version struct {
	XMLName xml.Name `xml:"version,omitempty" json:"version,omitempty"`
	Attrnumber string`xml:"number,attr"  json:",omitempty"`
}

