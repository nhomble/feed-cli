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
		Now: time.Time{},
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
	expected := "\n<!DOCTYPE html>\n<html>\n\t<body>\n\t\t<h1>0001-01-01 00:00:00 +0000 UTC</h1>\n\t\t\n\t\t<h1>someOrg</h1>\n\t\t<ul>\n\t\t\t\n\t\t\t<li><a href=\"someLink\">someArticle</a></li>\n\t\t\t\n\t\t</ul>\n\t\t\n\t</body>\n</html>\n"
	if expected != out {
		t.Fatalf("%s\n", out)
	}
}

func TestData_NowIn(t *testing.T) {
	utc, _ := time.LoadLocation("UTC")
	d := Data{
		Now: time.Date(2000, 1, 1, 1, 0, 0, 0, utc),
	}
	out := d.NowIn("EST")
	if "1999-12-31 20:00:00 -0500 EST" != out.String() {
		t.Fatalf("%s\n", out.String())
	}
}
