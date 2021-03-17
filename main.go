package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"x/bot"
	"x/config"
	"x/screenshot"
	"x/util"
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

	log.Printf("[Startup] Got info about %d politicians\n", poliStore.Len())

	client, user, err := bot.Login(cfg)
	if err != nil {
		panic("logging in to twitter: " + err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", user.ScreenName)

	events := wikipedia.StreamEdits(func(e *wikipedia.Event) bool {
		// If the article of an politician has been edited by an "anonymous" user (IP adress displayed)
		return poliStore.Contains(e.Title) && net.ParseIP(e.User) != nil
	})

	for edit := range events {
		fmt.Printf("%#v\n", edit)

		poli, ok := poliStore.Get(edit.Title)
		if !ok {
			continue
		}

		diffURL, ok := edit.DiffURL()
		if !ok {
			log.Printf("Couldn't generate diff URL for edit %#v\n", edit)
			continue
		}

		log.Println("Taking screenshot of", diffURL)
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

		var nameText string
		switch {
		case poli.FirstName == "" && poli.LastName != "":
			nameText = util.Hashtag(poli.LastName)
		case poli.FirstName != "" && poli.LastName != "":
			nameText = poli.FirstName + " " + util.Hashtag(poli.LastName)
		default:
			// data doesn't have a last name
			continue
		}

		var tweetText = fmt.Sprintf("Änderung beim Wiki-Eintrag von %s\n%s", nameText, diffURL)

		t, _, err := client.Statuses.Update(tweetText, &twitter.StatusUpdateParams{
			MediaIds: []int64{media.MediaID},
		})
		if err != nil {
			log.Printf("Error while tweeting: %s\n", err.Error())
			continue
		}

		fmt.Printf("Tweeted https://twitter.com/%s/status/%s\n", user.ScreenName, t.IDStr)
	}
}
