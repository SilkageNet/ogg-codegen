package tmpl

const schemaTmpl = `// Code generated by ogg-codegen. DO NOT EDIT.
package {{ opts.PackageName }}

import (
{{- range $key, $value := opts.Imports }}
"{{ $key }}"
{{- end }}
)

var (
{{- range $key, $value := opts.Imports }}
{{ $value }}
{{- end }}
)
{{range $type := .}}

{{ with $type.Schema.Comment }}{{ . }}{{ else }}// {{$type.Name}} defines model for {{$type.JSONName}}.{{ end }}
type {{$type.Name}} {{$type.Schema.TypeDecl}}

func (t *{{$type.Name}}) Valid() error {
var err error
{{$type.Schema.GenValidCode}}
return err
}

{{end}}
`