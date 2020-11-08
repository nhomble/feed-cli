package jobs

import (
	"bufio"
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/nhomble/feed-cli/template"
	"sync"
	"time"
)

var jobs chan Job
var feeds chan template.Entry

func setup(reader *bufio.Reader) {
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		if len(line) > 0 {
			jobs <- createJob(line)
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

	for job := range jobs {
		ctx, cancel := context.WithTimeout(context.Background(), job.timeout)
		defer cancel()
		feed, err := parser.ParseURLWithContext(job.link, ctx)
		if err == nil {
			if feed != nil {
				count := 0
				for _, item := range feed.Items {
					t := Item(*item).getTime().Add(job.age)
					if count > job.limit || t.Before(time.Now()) {
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
				fmt.Errorf("Did not parse a feed from line='%s'\n", job.link)
			}
		} else {
			fmt.Errorf("Failed to process='%s' err='%v'\n", job.link, err)
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

func Work(reader *bufio.Reader, workers int) chan template.Entry {
	jobs = make(chan Job)
	feeds = make(chan template.Entry)

	go setup(reader)
	go process(workers)

	return feeds
}
