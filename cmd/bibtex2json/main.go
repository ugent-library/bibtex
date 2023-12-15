package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ugent-library/bibtex"
)

func main() {
	p := bibtex.NewParser(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		e, err := p.Next()
		if err != nil {
			log.Fatal(err)
		}
		if e == nil {
			break
		}
		enc.Encode(e)
	}
}
