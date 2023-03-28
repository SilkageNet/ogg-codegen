package tmpl

const routerTmpl = `// Code generated by ogg-codegen. DO NOT EDIT.
package {{ opts.PackageName }}

import (
"github.com/gin-gonic/gin"
"go.uber.org/zap"
"net/http"
{{- range $key, $value := opts.Imports }}
"{{ $key }}"
{{- end }}
)

var (
{{- range $key, $value := opts.Imports }}
{{ $value }}
{{- end }}
)

{{ with .Comment }}{{ . }}{{ end }}

func (w *Wrapper) Setup{{ .GoName }}Routine() {
{{- range $op := .Operations}}
	w.engine.{{$op.Method}}("{{$op.Path}}", w.before("{{$op.OperationID}}"),{{with $op.NeedSecurity}} w.authVerifyFunc,{{end}} w.{{$op.OperationID}}, w.after("{{$op.OperationID}}"))
{{- end}}
}

{{- range $op := .Operations}}

{{ with $type := $op.ParamType }}
type {{$type.Name}} {{$type.Schema.TypeDecl}}

func (t *{{$type.Name}}) Bind(c *gin.Context) error {
{{$op.ReqParamBind}}
	return nil
}

func (t *{{$type.Name}}) Normalize() {
{{$op.ReqParamNormalize}}
}

func (t *{{$type.Name}}) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}
{{with $type.Schema.GenValidCode}}{{.}}{{end}}
	return nil
}
{{ end }}

{{ with $type := $op.BodyType }}
type {{$type.Name}} {{$type.Schema.TypeDecl}}

func (t *{{$type.Name}}) Bind(c *gin.Context) error {
{{$op.ReqBodyBind}}
}

func (t *{{$type.Name}}) Normalize() {
{{$op.ReqBodyNormalize}}
}

func (t *{{$type.Name}}) Valid() error {
	if t == nil {
		return errors.New("body is nil")
	}
{{with $type.Schema.GenValidCode}}{{.}}{{end}}
	return nil
}
{{ end }}

{{ with $op.CreateResponse }}
{{ with $type := $op.ResponseType }}
type {{$type.Name}} {{$op.ResponseTypeName}}
{{ end }}
{{ end }}

func (w *Wrapper) {{$op.OperationID}}(c *gin.Context) {
	var err error

{{- with $op.ParamType }}
	var param = &{{$op.ReqParamsName}}{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("{{$op.OperationID}}.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
{{- end }}

{{- with $op.ReqBodyName }}
	var body = &{{$op.ReqBodyName}}{}
	if err = w.bindAndValid(c, body); err != nil {
		w.logError("{{$op.OperationID}}.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
{{- end }}

{{- if $op.NoneLogic }}
	if err = w.server.{{$op.OperationID}}(c); err != nil {
		w.logError("{{$op.OperationID}}.err", zap.Error(err))
		w.errorWrapperFunc(c, err)
	}
{{- else }}
	w.logDebug("{{$op.OperationID}}"{{ with $op.ParamType }}, zap.Reflect("param", param){{ end }}{{ with $op.ReqBodyName }}, zap.Reflect("body", body){{ end }})
	var res {{$op.GetResponseName}}
	if res, err = w.server.{{$op.OperationID}}(c{{ with $op.ParamType }}, param{{ end }}{{ with $op.ReqBodyName}},body{{ end }}); err != nil {
		w.logError("{{$op.OperationID}}.err"{{ with $op.ParamType }}, zap.Reflect("param", param){{ end }}{{ with $op.ReqBodyName }}, zap.Reflect("body", body){{ end }}, zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("{{$op.OperationID}}.rsp"{{ with $op.ParamType }}, zap.Reflect("param", param){{ end }}{{ with $op.ReqBodyName }}, zap.Reflect("body", body){{ end }}{{ with $op.ServeFile }}{{ else }}, zap.Reflect("res", res){{ end }})
{{- end}}

	{{- if $op.ServeFile }}
	w.serveFile(c, res)
	{{- end }}
	{{- if $op.ServeText }}
	c.Data(http.StatusOK, "{{$op.ContentType}}", res)
	{{- end }}
	{{- if $op.ServeHTML }}
	c.Render(http.StatusOK, res)
	{{- end }}
	{{- if $op.ServeJSON }}
	c.JSON(http.StatusOK, res)
	{{- end }}
}

func (m *MockServer) {{$op.OperationID}}(c *gin.Context{{$op.ReqParamsNameWithParam}}{{$op.ReqBodyNameWithParam}}) {{$op.GetFuncReturn}} {
	return {{$op.GetEmptyReturn}}
}

{{- end}}
`
