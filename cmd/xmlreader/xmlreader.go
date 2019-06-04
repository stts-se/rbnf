package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stts-se/rbnf/xmlreader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "USAGE: xmlreader <xmlfiles>\n")
		os.Exit(1)
	}

	for _, f := range os.Args[1:] {
		_, err := xmlreader.RulesFromXMLFile(f)
		if err != nil {
			log.Fatalf("Couldn't parse file %s : %v", f, err)
		}
		log.Printf("Parsed xml rule file %s", f)
	}
}
