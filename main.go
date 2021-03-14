package main

import (
	"flag"
	"fmt"
	"log"
	"x/config"
	"x/twitter"
	"x/wikidata"
	"x/wikipedia"
)

var (
	flagConfigFile = flag.String("cfg", "config.yaml", "Config file path")
)

func main() {
	flag.Parse()

	cfg, err := config.Parse(*flagConfigFile)
	if err != nil {
		panic("parsing configuration file: " + err.Error())
	}

	log.Println("[Startup] Fetching politicians...")

	poliStore, err := wikidata.Politicians()
	if err != nil {
		panic("fetching politicians: " + err.Error())
	}

	client, user, err := twitter.Login(cfg)
	if err != nil {
		panic("logging in to twitter: " + err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", user.ScreenName)

	// Receive edit events and only return edits on sites of politicians
	events := wikipedia.StreamEdits(poliStore.Contains)

	for edit := range events {
		fmt.Printf("%#v\n", edit)

		// TODO
		t, _, err := client.Statuses.Update("", nil)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("Tweeted https://twitter.com/%s/%s\n", user.ScreenName, t.IDStr)
	}
}
