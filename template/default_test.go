package template

import (
	"bufio"
	"bytes"
	"testing"
	"time"
)

func TestDefaultTemplate(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	now := time.Now()
	Generate(writer, DefaultTemplateProvider{}, Data{
		Feeds: []Entry{
			{
				Article:   "someArticle",
				Link:      "someString",
				Org:       "someOrg",
				Published: &now,
			},
		},
	})
	writer.Flush()
	out := buf.String()
	expected := "\n<!DOCTYPE html>\n<html>\n\t<body>\n\t\t<ul>\n\t\t\t<li><a href=\"someString\">someArticle :: someOrg</a></li>\n\t\t</ul>\n\t</body>\n</html>\n"
	if expected != out {
		t.Fatalf("%s\n", out)
	}
}
