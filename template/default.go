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
		<ul>{{range .Feeds}}
			<li><a href="{{.Link}}">{{.Article}} :: {{.Org}}</a></li>
		{{end}}</ul>
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

type Entry struct {
	Article   string
	Link      string
	Org       string
	Published *time.Time
}

type Data struct {
	Feeds []Entry
}

func Generate(writer *bufio.Writer, provider Provider, data Data) {
	tpl := provider.GetTemplate()
	err := tpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}
