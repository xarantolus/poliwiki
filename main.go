package main

import (
	"flag"
	"fmt"
	"log"
	"x/config"
	"x/wikidata"
	"x/wikipedia"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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

	log.Println("Fetching politicians...")

	polis, err := wikidata.Politicians()
	if err != nil {
		panic("fetching politicians: " + err.Error())
	}

	var client *twitter.Client
	{
		config := oauth1.NewConfig(cfg.Twitter.APIKey, cfg.Twitter.APISecretKey)
		token := oauth1.NewToken(cfg.Twitter.AccessToken, cfg.Twitter.AccessTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)

		client = twitter.NewClient(httpClient)

		selfUser, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
		if err != nil {
			log.Fatalln("cannot log in with given credentials: " + err.Error())
		}

		log.Printf("[Twitter]: Logged in @%s\n", selfUser.ScreenName)
	}

	t, _, err := client.Statuses.Update("This is a test tweet", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(t.IDStr)

	log.Println("Start streaming edits...")

	events := wikipedia.StreamEdits(polis.Contains)

	for edit := range events {
		fmt.Printf("%#v\n", edit)
	}
}
