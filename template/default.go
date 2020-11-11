package template

import (
	"bufio"
	"text/template"
	"time"
)

type Provider interface {
	GetTemplate() *template.Template
}

type DefaultTemplateProvider struct {
}

func (p DefaultTemplateProvider) GetTemplate() *template.Template {
	t := `
<!DOCTYPE html>
<html>
	<body>
		<h1>{{with .Now}}{{.String}}{{end}}</h1>
		{{range .Feeds}}
		<h1>{{.Org}}</h1>
		<ul>
			{{range .Entries}}
			<li><a href="{{.Link}}">{{.Article}}</a></li>
			{{end}}
		</ul>
		{{end}}
	</body>
</html>
`
	tpl := template.New("default template")
	tpl, err := tpl.Parse(t)
	if err != nil {
		panic(err)
	}
	return tpl
}

type Feed struct {
	Org     string
	Entries []Entry
}

type Entry struct {
	Parent    Feed
	Article   string
	Link      string
	Published *time.Time
}

type Data struct {
	Feeds []Feed
	Now   time.Time
}

func Generate(writer *bufio.Writer, provider Provider, data Data) {
	tpl := provider.GetTemplate()
	err := tpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}
