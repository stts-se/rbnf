package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stts-se/rbnf/xmlreader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "USAGE: xmlreader <xmlfile> <input>\n")
		os.Exit(1)
	}

	f := os.Args[1]
	rPackage, err := xmlreader.RulesFromXMLFile(f)
	if err != nil {
		log.Fatalf("Couldn't parse file %s : %v", f, err)
	}
	log.Printf("Parsed xml rule file %s", f)

	if len(os.Args) > 2 {
		for _, s := range os.Args[2:] {
			res, err := rPackage.Spellout(s, "SpelloutRules", "spellout-numbering")
			if err != nil {
				log.Fatalf("Error for input string %s : %v", s, err)
			}
			fmt.Println(res)
		}
	}

}
