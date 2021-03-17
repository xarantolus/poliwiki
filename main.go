package main

import (
	"flag"
	"fmt"
	"log"
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
		return poliStore.Contains(e.Title)
	})

	// For detecting if we just tweeted about the same entry
	var (
		lastTweetID int64
		lastTitle   string
	)

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

		if edit.SizeDifference() < 50 {
			log.Println("Skipping small edit", diffURL)
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

		var tweetText = fmt.Sprintf("Änderung beim Wiki-Eintrag zu %s\n%s", nameText, diffURL)

		// If the last tweet was about this page, then we should just add it in a thread
		var replyID int64
		if lastTitle == edit.Title {
			replyID = lastTweetID
			tweetText = fmt.Sprintf("Noch eine Änderung bei %s\n%s", nameText, diffURL)
		}

		t, _, err := client.Statuses.Update(tweetText, &twitter.StatusUpdateParams{
			MediaIds:          []int64{media.MediaID},
			InReplyToStatusID: replyID,
		})
		if err != nil {
			log.Printf("Error while tweeting: %s\n", err.Error())
			continue
		}

		lastTitle = edit.Title
		lastTweetID = t.ID

		fmt.Printf("Tweeted https://twitter.com/%s/status/%s\n", user.ScreenName, t.IDStr)
	}
}
