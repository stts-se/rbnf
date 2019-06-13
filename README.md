# rbnf

Libraries and commands for numeric expansion, based on the "RBNF" part of the ICU project (http://site.icu-project.org/  and http://cldr.unicode.org/).

<!--
[![GoDoc](https://godoc.org/github.com/stts-se/rbnf?status.svg)](https://godoc.org/github.com/stts-se/rbnf) [![Go Report Card](https://goreportcard.com/badge/github.com/stts-se/rbnf)](https://goreportcard.com/report/github.com/stts-se/rbnf) [![Build Status](https://travis-ci.org/stts-se/rbnf.svg?branch=master)] 
-->

The license of the original software and data is here: https://github.com/unicode-org/icu/blob/master/icu4c/LICENSE http://www.unicode.org/copyright.html#License


The current spellout implementation does not use any of the original ICU code, but uses it supports most of the spellout rule format, and it can read the rule files, https://github.com/unicode-org/cldr/tree/master/common/rbnf.

## Unsupported formats 
The following formats are used in the ICU rules, but not supported by this package:
* _$_ (for singular and plural forms)
* _=#,##=_ and _=0=_
* _last primary ignorable_
* _→→→_


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
