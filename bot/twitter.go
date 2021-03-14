package bot

import (
	"x/config"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func Login(cfg config.Config) (client *twitter.Client, user *twitter.User, err error) {
	config := oauth1.NewConfig(cfg.Twitter.APIKey, cfg.Twitter.APISecretKey)
	token := oauth1.NewToken(cfg.Twitter.AccessToken, cfg.Twitter.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client = twitter.NewClient(httpClient)

	user, _, err = client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})

	return
}
