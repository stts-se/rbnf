#!/usr/bin/bash

set -e

### DEPENDENCIES #################
#
# 1) https://github.com/unicode-org/icu.git
#  $ git clone https://github.com/unicode-org/icu.git
#  $ cd icu/icu4j
#  $ ant
#  => icu4j.jar
#
# 2) github.com/HannaLindgren/go-utils/scripts/compare_line_by_line/compare_line_by_line.go
#  $ git clone https://github.com/HannaLindgren/go-utils/
#  $ cd go-utils/scripts/compare_line_by_line
#  $ go build
#  => compare_line_by_line
#
##################################


icu4j_jar=~/progz/icu/icu4j/icu4j.jar
outdir="output-ordinal"
#numsfile=nums_1_to_100k
numsfile=nums_1_to_1M

mkdir -p $outdir


all_langs="af ak am ar az be ca ccp da de de_CH ee el en en_IN eo es ff fi fil fr fr_BE fr_CH he hi hr hu id it ja km ko ky lb lo ms my nb nl pt pt_PT qu ru sr sr_Latn sv sw ta th tr vi yue yue_Hans zh"

declare -A rulesets=( ["ar"]="spellout-ordinal-feminine" ["be"]="spellout-ordinal-feminine" ["ca"]="spellout-ordinal-feminine" ["da"]="spellout-ordinal-common" ["el"]="spellout-ordinal-feminine" ["es"]="spellout-ordinal-feminine" ["fr"]="spellout-ordinal-feminine" ["fr_BE"]="spellout-ordinal-feminine" ["fr_CH"]="spellout-ordinal-feminine" ["he"]="spellout-ordinal-feminine" ["hi"]="spellout-ordinal-feminine" ["hr"]="spellout-ordinal-feminine" ["it"]="spellout-ordinal-feminine" ["ko"]="spellout-ordinal-native" ["lb"]="spellout-ordinal-feminine" ["nb"]="spellout-ordinal-feminine" ["pt"]="spellout-ordinal-feminine" ["pt_PT"]="spellout-ordinal-feminine" ["ru"]="spellout-ordinal-feminine" ["sv"]="spellout-ordinal-feminine" )

langs=$all_langs
#langs="sv en de fr da fi es"

echo "LANGS: $langs" 1>&2
echo "OUTDIR: $outdir" 1>&2
echo "NUMSFILE: $numsfile" 1>&2
echo "" 1>&2

for lang in $langs; do

    echo "=== PROCESSING $lang ... " 1>&2
    outicu4j="$outdir/${lang}_${numsfile}_icu4j.txt"
    outgo="$outdir/${lang}_${numsfile}_rbnfgo.txt"
    outdiff="$outdir/${lang}_${numsfile}.diff"

    # rule set name
    ruleset=${rulesets[$lang]}
    if [ "<$ruleset>" == "<>" ]; then
	ruleset="spellout-ordinal"
    fi
    echo "Rule set: $ruleset" 1>&2

    # list ordinal rule sets:
    # go run ../cmd/spellout/spellout.go -l https://github.com/unicode-org/cldr/raw/master/common/rbnf/${lang}.xml | egrep "spellout-ordinal" | head -1 | sed "s/^/${lang}\t/" | sed 's/ - //' | sed 's/ .public.*//' | egrep -v "spellout-ordinal$" 
    
    time cat ${numsfile}.txt | scala -cp ${icu4j_jar} batch_run_icu4j_ordinal.scala ${lang} ${ruleset} >| $outicu4j
    if time cat ${numsfile}.txt | go run ../cmd/spellout/spellout.go -r $ruleset https://github.com/unicode-org/cldr/raw/master/common/rbnf/${lang}.xml >| $outgo ; then
    	compare_line_by_line -q $outgo $outicu4j &> $outdiff
    	echo "=== $lang DONE" 1>&2
    else
    	echo "=== $lang FAILED" 1>&2
    fi
    echo "" 1>&2
	

done
