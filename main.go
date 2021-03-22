package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xarantolus/poliwiki/bot"
	"github.com/xarantolus/poliwiki/config"
	"github.com/xarantolus/poliwiki/screenshot"
	"github.com/xarantolus/poliwiki/util"
	"github.com/xarantolus/poliwiki/wikidata"
	"github.com/xarantolus/poliwiki/wikipedia"

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

	type lastInfo struct {
		TweetID int64
		Time    time.Time
	}

	// For detecting if we tweeted about the same entry within the last two hours
	var lastTweetInfo = make(map[string]lastInfo)

	for edit := range events {
		log.Printf("%#v\n", edit)

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
		case poli.Name != "":
			nameText = poli.Name
		default:
			// data doesn't have a name, shouldn't really happen?
			continue
		}

		var partyText = poli.PartyShortname()
		if partyText != "" {
			partyText = ", " + util.Hashtag(partyText)
		}

		var tweetText = fmt.Sprintf("Änderung beim Wiki-Eintrag zu %s%s\n%s", nameText, partyText, diffURL)

		// If we tweeted about this in the last two hours, add it in a thread
		var replyID int64
		if li := lastTweetInfo[edit.Title]; time.Since(li.Time) < 2*time.Hour {
			replyID = li.TweetID
			tweetText = fmt.Sprintf("Noch eine Änderung bei %s%s\n%s", nameText, partyText, diffURL)
		}

		t, _, err := client.Statuses.Update(tweetText, &twitter.StatusUpdateParams{
			MediaIds:          []int64{media.MediaID},
			InReplyToStatusID: replyID,
		})
		if err != nil {
			log.Printf("Error while tweeting: %s\n", err.Error())
			continue
		}

		// Save this info for the next tweet
		lastTweetInfo[edit.Title] = lastInfo{
			TweetID: t.ID,
			Time:    time.Now(),
		}

		log.Printf("Tweeted https://twitter.com/%s/status/%s\n", user.ScreenName, t.IDStr)
	}
}
