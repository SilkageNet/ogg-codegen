package tmpl

const serverTmpl = `// Code generated by ogg-codegen. DO NOT EDIT.
package {{ opts.PackageName }}

import (
"github.com/gin-gonic/gin"
"github.com/gin-gonic/gin/render"
)

var _ render.Render

type Server interface {
{{- range $tag := .}}
{{$tag.Comment}}
{{- range $op:= $tag.Operations}}
	{{$op.Comment}}
	{{$op.OperationID}}(c *gin.Context{{$op.ReqParamsNameWithParam}}{{$op.ReqBodyNameWithParam}}) {{$op.GetFuncReturn}}
{{- end}}
{{end}}
}

func Register(engine *gin.Engine, server Server, fns ...OptionFunc) {
	var opt = defaultOption
	for _, fn := range fns {
		fn(opt)
	}
	var wrapper = &Wrapper{option: *opt, engine: engine, server: server}
{{- range $tag := .}}
	wrapper.Setup{{$tag.GoName}}Routine()
{{- end}}
}

type MockServer struct {
}
`
