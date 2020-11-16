package main

import (
	"bufio"
	"flag"
	"github.com/nhomble/feed-cli/jobs"
	"github.com/nhomble/feed-cli/template"
	"os"
	"time"
)

var (
	workers          = flag.Int("numWorkers", 5, "Number of worker threads")
	templateOverride = flag.String("templateOverride", "", "Relative path to template override")
)

func main() {
	flag.Parse()

	if *workers < 1 {
		panic("Need to have a positive number of workers!")
	}

	var provider template.Provider
	if len(*templateOverride) == 0 {
		provider = template.DefaultTemplateProvider{}
	} else {
		provider = template.RemoteTemplateProvider{Override: *templateOverride}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	feeds := jobs.Work(bufio.NewReader(os.Stdin), *workers)

	template.Generate(writer, provider, template.Data{
		Feeds: feeds,
		Now:   time.Now(),
	})
}
