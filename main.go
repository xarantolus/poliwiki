package main

import (
	"flag"
	"fmt"
	"log"
	"x/bot"
	"x/config"
	"x/screenshot"
	"x/wikidata"
	"x/wikipedia"

	"github.com/dghubble/go-twitter/twitter"
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

	client, user, err := bot.Login(cfg)
	if err != nil {
		panic("logging in to twitter: " + err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", user.ScreenName)

	// Receive edit events and only return edits on sites of politicians
	events := wikipedia.StreamEdits(poliStore.Contains)

	for edit := range events {
		fmt.Printf("%#v\n", edit)

		diffURL, ok := edit.DiffURL()
		if !ok {
			log.Printf("Couldn't generate diff URL for edit %#v\n", edit)
			continue
		}

		png, err := screenshot.Take(diffURL)
		if err != nil {
			log.Printf("Error while taking screenshot: %s\n", err.Error())
			continue
		}

		media, _, err := client.Media.Upload(png, "image/png")
		if err != nil {
			log.Printf("Error while uploading image: %s\n", err.Error())
			continue
		}

		// TODO: generate text
		t, _, err := client.Statuses.Update("", &twitter.StatusUpdateParams{
			MediaIds: []int64{media.MediaID},
		})
		if err != nil {
			log.Printf("Error while tweeting: %s\n", err.Error())
			continue
		}

		fmt.Printf("Tweeted https://twitter.com/%s/%s\n", user.ScreenName, t.IDStr)
	}
}
