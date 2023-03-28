package tmpl

import (
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"strings"
	"text/template"
)

var templateFuncMap = template.FuncMap{
	"lower":      strings.ToLower,
	"title":      strings.Title,
	"lcFirst":    util.LowercaseFirstCharacter,
	"ucFirst":    util.UppercaseFirstCharacter,
	"camelCase":  util.ToCamelCase,
	"tmplSymbol": func() string { return "`" },
}

type Options struct {
	PackageName string
	Imports     map[string]string
}

func RegisterFuncMap(key string, fn interface{}) {
	templateFuncMap[key] = fn
}

func RegisterOptionsFunc(opts *Options) {
	templateFuncMap["opts"] = func() Options { return *opts }
}
