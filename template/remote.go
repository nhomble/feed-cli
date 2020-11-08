package template

import "text/template"

type RemoteTemplateProvider struct {
	Override string
}

func (p RemoteTemplateProvider) GetTemplate() *template.Template {
	tpl, err := template.ParseFiles(p.Override)
	if err != nil {
		panic(err)
	}
	return tpl
}
