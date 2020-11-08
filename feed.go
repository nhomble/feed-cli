package main

import (
	"bufio"
	"fmt"
	"github.com/mmcdole/gofeed"
	"os"
	"sync"
	"text/template"
	"time"
)

const WORKERS = 5
const LIMIT = 5
const DAYS = 24 * time.Hour

const defaultTemplate = `
<html>
	<body>
		<ul>{{range .Feeds}}
			<li><a href="{{.Link}}">{{.Article}} :: {{.Org}}</a></li>
		{{end}}</ul>
	</body>
</html>
`

type Entry struct {
	Article   string
	Link      string
	Org       string
	published *time.Time
}

type Data struct {
	Feeds chan Entry
}

func Generate(writer *bufio.Writer, data Data) {
	tpl := template.New("feed template")
	tpl, err := tpl.Parse(defaultTemplate)
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}

var jobs chan string
var feeds chan Entry

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
					feeds <- Entry{
						Article:   item.Title,
						Link:      item.Link,
						Org:       feed.Title,
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

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	jobs = make(chan string)
	feeds = make(chan Entry)

	go setup()
	go process(WORKERS)
	Generate(writer, Data{
		Feeds: feeds,
	})
}
