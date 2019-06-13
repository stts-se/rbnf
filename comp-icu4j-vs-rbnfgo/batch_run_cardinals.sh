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

mkdir -p output

icu4j_jar=~/progz/icu/icu4j/icu4j.jar

#numsfile=nums_1_to_100k
numsfile=nums_1_to_1M

all_langs="af ak am ar az be bg bs ca ccp chr cs cy da de de_CH ee el en en_001 en_IN eo es es_419 et fa fa_AF ff fi fil fo fr fr_BE fr_CH ga he hi hr hu hy id is it ja ka kl km ko ky lb lo lrc lt lv mk ms mt my nb nl nn pl pt pt_PT qu ro root ru se sk sl sq sr sr_Latn sv sw ta th tr uk vi yue yue_Hans zh zh_Hant"
tested_langs="en sv fr de da ja ar cs fi en_IN hu it ta"
# not working: th (not sure why); ru (sin/plu); sk (sin/plu)

langs=$all_langs

for lang in $langs; do

    echo "=== PROCESSING $lang ... " 1>&2
    outgo="output/${numsfile}_out_${lang}_rbnfgo.txt"
    outicu4j="output/${numsfile}_out_${lang}_icu4j.txt"

    # rule set name
    ruleset="spellout-numbering"
    if [ $lang == "da" ]; then
	ruleset="spellout-cardinal-common"
    fi
    echo "Rule set: $ruleset" 1>&2

    
    time cat ${numsfile}.txt | scala -cp ${icu4j_jar} batch_run_icu4j.scala ${lang} >| $outicu4j
    time cat ${numsfile}.txt | go run ../cmd/spellout/spellout.go -r $ruleset https://github.com/unicode-org/cldr/raw/master/common/rbnf/${lang}.xml >| $outgo

    compare_line_by_line -q $outgo $outicu4j &> output/${lang}.diff
    echo "=== DONE" 1>&2
    echo "" 1>&2

done;
