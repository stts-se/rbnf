# rbnf

This is a Go implementation of parts of ICU's spellout rules (RBNF) for numbers.


It's based on the "RBNF" part of the ICU project (http://site.icu-project.org/  and http://cldr.unicode.org/).

<!--
[![GoDoc](https://godoc.org/github.com/stts-se/rbnf?status.svg)](https://godoc.org/github.com/stts-se/rbnf) [![Go Report Card](https://goreportcard.com/badge/github.com/stts-se/rbnf)](https://goreportcard.com/report/github.com/stts-se/rbnf) [![Build Status](https://travis-ci.org/stts-se/rbnf.svg?branch=master)] 
-->


The license of the original software and data is here: https://github.com/unicode-org/icu/blob/master/icu4c/LICENSE http://www.unicode.org/copyright.html#License


The current spellout implementation does not use any of the original ICU code, but uses it supports most of the spellout rule format, and it can read the rule files, https://github.com/unicode-org/cldr/tree/master/common/rbnf.

## Unsupported features
The following format strings are used in the ICU rules, but not supported by this package:
* _$_ (for singular and plural forms)
* _=#,##=_ and _=0=_
* _last primary ignorable_
* _→→→_

The rule sets have information on the public/private attribute, but the distinction is not supported on rule execution (all rules can be references as public).


## Command line tool

There is a command line program for expanding numbers according to a CLDR rule file.

Here's an example, using the English rules.

You need to have Go >= 1.12 installed.

Start by cloning the `github.com/stts-se/rbnf` repository, then:

    cd rbnf/cmd/spellout/
    go build
    

List the rules of the English rule file (reading the file directly from github):

     ./spellout -l https://raw.githubusercontent.com/unicode-org/cldr/master/common/rbnf/en.xml
     == Listing public rule sets ==
     SpelloutRules
      - spellout-cardinal [public] (38)
      - spellout-cardinal-verbose [public] (12)
      - spellout-numbering [public] (4)
      - spellout-numbering-verbose [public] (4)
      - spellout-numbering-year [public] (29)
      - spellout-ordinal [public] (30)
      - spellout-ordinal-verbose [public] (10)
      OrdinalRules
       - digits-ordinal [public] (1)



Or, download the rule file and save as en.xml:

       curl https://raw.githubusercontent.com/unicode-org/cldr/master/common/rbnf/en.xml > en.xml
    

Test the default cardinal rule expansion:

      ./spellout -r spellout-numbering en.xml 1066
      1066	one thousand sixty-six

Test spelling out as year:

      ./spellout -r spellout-numbering-year en.xml 1066
      1066	ten sixty-six












## Links and references
* ICU project and source code: <br/>
  http://site.icu-project.org/ <br/>
  https://github.com/unicode-org/icu
* ICU user guide: http://userguide.icu-project.org/
* CLDR project and source code: <br/>
  http://cldr.unicode.org/ <br/>
  https://github.com/unicode-org/cldr <br/>
* Rule format: <br/>
  http://userguide.icu-project.org/formatparse/numbers/rbnf-examples <br/>
  http://icu-project.org/apiref/icu4c/classRuleBasedNumberFormat.html   <br/>
  http://www.icu-project.org/applets/icu4j/4.1/docs-4_1_1/com/ibm/icu/text/RuleBasedNumberFormat.html 
* Rule files: https://github.com/unicode-org/cldr/tree/master/common/rbnf
* ICU license: <br/>
  http://www.unicode.org/copyright.html#License <br/>
  https://github.com/unicode-org/icu/blob/master/icu4c/LICENSE
* Go implementation of parts of CLDR: <br/>
  https://godoc.org/golang.org/x/text/unicode/cldr <br/>
  https://github.com/golang/text/tree/master/unicode/cldr 
