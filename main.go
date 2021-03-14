package main

import (
	"fmt"
	"x/wikidata"
	"x/wikipedia"
)

func main() {
	polis, err := wikidata.Politicians()
	if err != nil {
		panic(err)
	}

	events := wikipedia.StreamEdits(polis.Contains)

	for edit := range events {
		fmt.Printf("%#v\n", edit)
	}
}
