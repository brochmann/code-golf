package oauth

import (
	"net/url"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
	"golang.org/x/oauth2/stackoverflow"
)

type Config struct {
	oauth2.Config
	Name, UserEndpoint string
}

var Connections = map[string]*Config{
	// https://discord.com/developers/applications
	"discord": {
		Name:         "Discord",
		UserEndpoint: discordgo.EndpointUser("@me"),
		Config: oauth2.Config{
			Scopes: []string{"identify"},
			Endpoint: oauth2.Endpoint{
				AuthStyle: oauth2.AuthStyleInParams,
				AuthURL:   discordgo.EndpointOauth2 + "authorize",
				TokenURL:  discordgo.EndpointOauth2 + "token",
			},
		},
	},

	// https://github.com/settings/developers
	"github": {
		Name:   "GitHub",
		Config: oauth2.Config{Endpoint: github.Endpoint},
	},

	// https://gitlab.com/-/profile/applications
	"gitlab": {
		Name:         "GitLab",
		UserEndpoint: "https://gitlab.com/oauth/userinfo",
		Config: oauth2.Config{
			Endpoint: gitlab.Endpoint,
			Scopes:   []string{"openid"},
		},
	},

	// https://stackapps.com/apps/oauth
	"stack-overflow": {
		Name:         "Stack Overflow",
		Config:       oauth2.Config{Endpoint: stackoverflow.Endpoint},
		UserEndpoint: "https://api.stackexchange.com/me?site=stackoverflow",
	},
}

func init() {
	host := "code.golf"
	if _, dev := os.LookupEnv("DEV"); dev {
		host = "localhost"
	}

	for id, config := range Connections {
		prefix := strings.ReplaceAll(strings.ToUpper(id), "-", "_")

		config.ClientID = os.Getenv(prefix + "_CLIENT_ID")
		config.ClientSecret = os.Getenv(prefix + "_CLIENT_SECRET")
		config.RedirectURL = "https://" + host + "/golfer/connect/" + id

		// Add a key to UserEndpoint if we have one.
		if key := os.Getenv(prefix + "_KEY"); key != "" {
			u, err := url.Parse(config.UserEndpoint)
			if err != nil {
				panic(err)
			}

			q := u.Query()
			q.Set("key", key)
			u.RawQuery = q.Encode()
			config.UserEndpoint = u.String()
		}
	}
}
