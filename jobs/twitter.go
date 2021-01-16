package jobs

import (
	"context"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/nhomble/feed-cli/template"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"strings"
	"time"
)

const twitterScheme = "twitter://@"

func (j Job) isTwitter() bool {
	return strings.HasPrefix(j.link, twitterScheme)
}

func parseTwitterUser(job Job) string {
	return strings.Replace(job.link, twitterScheme, "", -1)
}

type twitterParser struct {
	client *twitter.Client
}

func (p twitterParser) parse(job Job) []template.Feed {
	username := parseTwitterUser(job)
	ret := make([]template.Feed, 0)
	tweets, _, _ := p.client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: parseTwitterUser(job),
		Count:      job.limit,
	})
	if tweets != nil {
		group := template.Feed{
			Org:     fmt.Sprintf("@%s", username),
			Entries: []template.Entry{},
		}
		for _, tweet := range tweets {
			t, _ := tweet.CreatedAtTime()
			if t.Add(job.age).Before(time.Now()) {
				continue
			}
			group.Entries = append(group.Entries, template.Entry{
				Article:   tweet.Text,
				Link:      fmt.Sprintf("https://mobile.twitter.com/%s/status/%s", username, tweet.IDStr),
				Published: &t,
			})
		}
		if len(group.Entries) > 0 {
			ret = append(ret, group)
		}
	}
	return ret
}

func createTwitterParser() twitterParser {
	key := os.Getenv("TWITTER_CONSUMER_KEY")
	secret := os.Getenv("TWITTER_CONSUMER_SECRET")
	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(context.Background())

	client := twitter.NewClient(httpClient)
	return twitterParser{client: client}
}
