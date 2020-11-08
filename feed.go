package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/nhomble/feed-cli/template"
	"os"
	"sync"
	"time"
)

const WORKERS = 5
const LIMIT = 5
const DAYS = 24 * time.Hour
const TIMEOUT = 1 * time.Minute

var jobs chan string
var feeds chan template.Entry

func setup() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		if len(line) > 0 {
			jobs <- string(line)
		}
	}
	close(jobs)
}

func process(workers int) {
	var group sync.WaitGroup

	for i := 0; i < workers; i++ {
		group.Add(1)
		go worker(&group)
	}
	group.Wait()
	close(feeds)
}

func worker(group *sync.WaitGroup) {
	parser := gofeed.NewParser()
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	for line := range jobs {
		feed, err := parser.ParseURLWithContext(line, ctx)
		if err == nil {
			if feed != nil {
				count := 0
				for _, item := range feed.Items {
					t := Item(*item).getTime().Add(14 * DAYS)
					if count > LIMIT || t.Before(time.Now()) {
						break
					}
					count++
					feeds <- template.Entry{
						Article:   item.Title,
						Link:      item.Link,
						Org:       feed.Title,
						Published: item.PublishedParsed,
					}
				}
			} else {
				fmt.Errorf("Did not parse a feed from line='%s'\n", line)
			}
		} else {
			fmt.Errorf("Failed to process='%s' err='%v'\n", line, err)
		}
	}
	group.Done()
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

func main() {
	var templateOverride string
	flag.StringVar(&templateOverride, "templateOverride", "", "Relative path to template override")
	flag.Parse()

	var provider template.TemplateProvider
	if len(templateOverride) == 0 {
		provider = template.DefaultTemplateProvider{}
	} else {
		provider = template.RemoteTemplateProvider{templateOverride}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	jobs = make(chan string)
	feeds = make(chan template.Entry)

	go setup()
	go process(WORKERS)
	template.Generate(writer, provider, template.Data{
		Feeds: feeds,
	})
}
