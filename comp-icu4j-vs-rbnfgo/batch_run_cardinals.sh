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
outdir="output-cardinal"
#numsfile=nums_1_to_100k
numsfile=nums_1_to_1M

mkdir -p $outdir


all_langs="af ak am ar az be bg bs ca ccp chr cs cy da de de_CH ee el en en_IN eo es et fa fa_AF ff fi fil fo fr fr_BE fr_CH ga he hi hr hu hy id is it ja ka kl km ko ky lb lo lrc lt lv mk ms mt my nb nl nn pl pt pt_PT qu ro ru se sk sl sq sr sr_Latn sv sw ta th tr uk vi yue yue_Hans zh"

langs=$all_langs

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
    ruleset="spellout-numbering"
    if [ $lang == "da" ]; then
	ruleset="spellout-cardinal-common"
    fi
    echo "Rule set: $ruleset" 1>&2

    
    time cat ${numsfile}.txt | scala -cp ${icu4j_jar} batch_run_icu4j.scala ${lang} >| $outicu4j
    if time cat ${numsfile}.txt | go run ../cmd/spellout/spellout.go -r $ruleset https://github.com/unicode-org/cldr/raw/master/common/rbnf/${lang}.xml >| $outgo ; then
	compare_line_by_line -q $outgo $outicu4j &> $outdiff
	echo "=== $lang DONE" 1>&2
    else
	echo "=== $lang FAILED" 1>&2
    fi
    echo "" 1>&2
	

done
