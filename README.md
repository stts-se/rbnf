# rbnf

Libraries and commands for spelling out numeric expansion, based on the "RBNF" the ICU project, https://github.com/unicode-org/icu and https://github.com/unicode-org/cldr.

The license of the original software and data is here: https://github.com/unicode-org/icu/blob/master/icu4c/LICENSE


The current spellout implementation does not use any of the original ICU code, but uses the spellout rule format and data files, https://github.com/unicode-org/cldr/tree/master/common/rbnf.

Unsupported formats from the ICU rules:
* _$_ (for singular and plural forms)
* _=#,##=_ and _=0=_
* _last primary ignorable_
* _→→→_


Links:
* ICU project: http://userguide.icu-project.org/
* Rule format: http://icu-project.org/apiref/icu4c/classRuleBasedNumberFormat.html | http://www.icu-project.org/applets/icu4j/4.1/docs-4_1_1/com/ibm/icu/text/RuleBasedNumberFormat.html
* Rule files (xml): https://github.com/unicode-org/cldr/tree/master/common/rbnf
* Go implementation of parts of CLDR: https://godoc.org/golang.org/x/text/unicode/cldr and https://github.com/golang/text/tree/master/unicode/cldr 
* ICU license: https://github.com/unicode-org/icu/blob/master/icu4c/LICENSE
