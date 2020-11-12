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
var feeds chan template.Feed

func setup(reader *bufio.Reader) {
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		if canProcess(line) {
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
					feeds <- group
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

func Work(reader *bufio.Reader, workers int) []template.Feed {
	jobs = make(chan Job)
	feeds = make(chan template.Feed)

	go setup(reader)
	go process(workers)

	ret := []template.Feed{}
	for val := range feeds {
		ret = append(ret, val)
	}
	return ret
}
