package tmpl

import (
	"text/template"
)

const (
	tmplName = "ogg-codegen"
)

var templates = map[string]string{
	"schema.tmpl":  schemaTmpl,
	"enum.tmpl":    enumTmpl,
	"wrapper.tmpl": wrapperTmpl,
	"server.tmpl":  serverTmpl,
	"router.tmpl":  routerTmpl,
}

func Parse() (*template.Template, error) {
	var t = template.New(tmplName).Funcs(templateFuncMap)
	for name, s := range templates {
		var tmpl = t.New(name)
		var _, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
