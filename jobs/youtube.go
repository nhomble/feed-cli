package jobs

import (
	"fmt"
	soup2 "github.com/anaskhan96/soup"
	"github.com/nhomble/feed-cli/template"
	"strings"
)

const youtubeScheme = "youtube_user://@"

type youtubeParser struct {
	rssParser goFeedParser
}

func createYoutubeParser() youtubeParser {
	return youtubeParser{
		rssParser: defaultParser,
	}
}

func parseYoutubeUser(job Job) string {
	return strings.Replace(job.link, youtubeScheme, "", -1)
}

func (j Job) isYoutube() bool {
	return strings.HasPrefix(j.link, youtubeScheme)
}

func (p youtubeParser) parse(job Job) []template.Feed {
	// get rss url from user
	youtubeUser := parseYoutubeUser(job)
	youtubeLink := fmt.Sprintf("https://youtube.com/@%s", youtubeUser)

	resp, err := soup2.Get(youtubeLink)
	if err != nil {
		return make([]template.Feed, 0)
	}
	doc := soup2.HTMLParse(resp)
	links := doc.FindAll("link")
	for _, link := range links {
		if val, ok := link.Attrs()["title"]; ok && val == "RSS" {
			rss, ok := link.Attrs()["href"]
			if ok {
				return p.rssParser.parse(Job{
					link:         rss,
					limit:        job.limit,
					age:          job.age,
					timeout:      job.timeout,
					nameOverride: job.nameOverride,
				})
			}
		}

	}

	return make([]template.Feed, 0)
}
