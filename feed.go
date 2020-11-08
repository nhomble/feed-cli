package main

import (
	"bufio"
	"flag"
	"github.com/nhomble/feed-cli/jobs"
	"github.com/nhomble/feed-cli/template"
	"os"
)

func main() {
	var templateOverride string
	flag.StringVar(&templateOverride, "templateOverride", "", "Relative path to template override")
	flag.Parse()

	var provider template.TemplateProvider
	if len(templateOverride) == 0 {
		provider = template.DefaultTemplateProvider{}
	} else {
		provider = template.RemoteTemplateProvider{Override: templateOverride}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	feeds := jobs.Work(bufio.NewReader(os.Stdin))

	template.Generate(writer, provider, template.Data{
		Feeds: feeds,
	})
}
