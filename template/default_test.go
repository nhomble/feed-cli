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
		Feeds: []Feed{
			{
				Org: "someOrg",
				Entries: []Entry{
					{
						Article:   "someArticle",
						Link:      "someLink",
						Published: &now,
					},
				},
			},
		},
	})
	writer.Flush()
	out := buf.String()
	expected := "\n<!DOCTYPE html>\n<html>\n\t<body>\n\t\t\n\t\t<h1>someOrg</h1>\n\t\t<ul>\n\t\t\t\n\t\t\t<li><a href=\"someLink\">someArticle</a></li>\n\t\t\t\n\t\t</ul>\n\t\t\n\t</body>\n</html>\n"
	if expected != out {
		t.Fatalf("%s\n", out)
	}
}
