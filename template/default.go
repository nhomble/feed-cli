package template

import "text/template"

type TemplateProvider interface {
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
