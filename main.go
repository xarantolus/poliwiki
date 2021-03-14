package main

import (
	"fmt"
	"x/wikidata"
)

func main() {
	polis, err := wikidata.Politicians()
	if err != nil {
		panic(err)
	}
	fmt.Println(polis.Get("Angela Merkel"))
}
