package main

import (
	"bufio"
	"fmt"
	"github.com/mmcdole/gofeed"
	"os"
	"sync"
	"time"
)

const WORKERS = 5
const LIMIT = 5
const DAYS = 24 * time.Hour

var jobs chan string
var feeds chan entry

type entry struct {
	article   string
	link      string
	org       string
	published *time.Time
}

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

	for line := range jobs {
		feed, err := parser.ParseURL(line)
		if err == nil {
			if feed != nil {
				count := 0
				for _, item := range feed.Items {
					t := Item(*item).getTime().Add(14 * DAYS)
					if count > LIMIT || t.Before(time.Now()) {
						break
					}
					count++
					feeds <- entry{
						article:   item.Title,
						link:      item.Link,
						org:       feed.Title,
						published: item.PublishedParsed,
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

func write(writer *bufio.Writer) {
	for entry := range feeds {
		writer.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s :: %s</a></li>\n", entry.link, entry.article, entry.org))
	}
}

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	jobs = make(chan string)
	feeds = make(chan entry)

	writer.WriteString("<html>\n")
	writer.WriteString("<body>\n")
	writer.WriteString("<ul>\n")

	go setup()
	go process(WORKERS)
	write(writer)

	writer.WriteString("</ul>\n")
	writer.WriteString("</body>\n")
	writer.WriteString("</html>")
}
