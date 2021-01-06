package jobs

import (
	"bufio"
	"github.com/nhomble/feed-cli/template"
	"sync"
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

	for job := range jobs {
		for _, group := range parserForJob(job).parse(job) {
			feeds <- group
		}
	}
	group.Done()
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
