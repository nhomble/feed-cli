package jobs

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/nhomble/feed-cli/template"
	"time"
)

var defaultParser = goFeedParser{
	gofeed.NewParser(),
}

type FeedParser interface {
	parse(Job) []template.Feed
}

type goFeedParser struct {
	parser *gofeed.Parser
}

func (parser goFeedParser) parse(job Job) []template.Feed {
	ctx, cancel := context.WithTimeout(context.Background(), job.timeout)
	defer cancel()
	feed, err := parser.parser.ParseURLWithContext(job.link, ctx)
	if err != nil {
		fmt.Errorf("Failed to process='%s' err='%v'\n", job.link, err)
		return []template.Feed{}
	}
	return transform(job, feed)
}

func parserForJob(job Job) FeedParser {
	return defaultParser
}

type Item gofeed.Item

func (i Item) getTime() time.Time {
	if i.UpdatedParsed != nil {
		return *i.UpdatedParsed
	}
	if i.PublishedParsed != nil {
		return *i.PublishedParsed
	}
	return time.Now()
}

func transform(job Job, feed *gofeed.Feed) []template.Feed {
	ret := make([]template.Feed, 0)
	if feed != nil {
		count := 1
		group := template.Feed{
			Org:     feed.Title,
			Entries: []template.Entry{},
		}
		for _, item := range feed.Items {
			t := Item(*item).getTime().Add(job.age)
			if count > job.limit || t.Before(time.Now()) {
				break
			}
			count++
			group.Entries = append(group.Entries, template.Entry{
				Article:   item.Title,
				Link:      item.Link,
				Published: item.PublishedParsed,
			})
		}
		if len(group.Entries) > 0 {
			ret = append(ret, group)
		}
	}
	return ret
}
