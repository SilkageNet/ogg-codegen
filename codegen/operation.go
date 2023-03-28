package codegen

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/getkin/kin-openapi/openapi3"
	"net/http"
	"strings"
)

type Operation struct {
	OperationID string              // 唯一表示标识
	Path        string              // 路径
	Method      string              // 请求类型
	Comment     string              // 注释
	Spec        *openapi3.Operation // OpenAPIv3 original Operation
	ServeFile   bool
	ServeText   bool
	ServeHTML   bool
	ServeJSON   bool

	NoneLogic bool

	HeaderParams []Parameter
	QueryParams  []Parameter
	CookieParams []Parameter
	PathParams   []Parameter
	Body         *RequestBody
	Securities   []Security

	ParamType    *Type
	BodyType     *Type
	ResponseType *Type
	ContentType  string
}

func (o Operation) AllParams() []Parameter {
	var params = o.QueryParams
	params = append(params, o.HeaderParams...)
	params = append(params, o.CookieParams...)
	params = append(params, o.PathParams...)
	return params
}

func (o Operation) ReqParamsName() string {
	return fmt.Sprintf("%sParams", o.OperationID)
}

func (o Operation) ReqParamsNameWithParam() string {
	if len(o.AllParams()) == 0 {
		return ""
	}
	return fmt.Sprintf(", params *%s", o.ReqParamsName())
}

func (o Operation) ReqParamBind() string {
	if o.ParamType == nil {
		return ""
	}
	var parts []string
	for _, p := range o.QueryParams {
		var part string
		switch {
		case util.IsGoStrType(p.Schema.GoType):
			part = fmt.Sprintf(`t.%s = c.Query("%s")`, p.GoFieldName(), p.Name)
		case util.IsGoNumType(p.Schema.GoType):
			part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.Query("%s")))`, p.GoFieldName(), p.Schema.GoType, p.Name)
		default:
			var t = p.Schema.RefType
			if t == "" {
				t = p.Schema.GoType
			}
			switch {
			case t == "tex.JsInt64" || util.IsGoNumType(p.Schema.Spec.Format):
				part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.Query("%s")))`, p.GoFieldName(), t, p.Name)
			case util.IsGoStrType(p.Schema.Spec.Type):
				part = fmt.Sprintf(`t.%s = %s(c.Query("%s"))`, p.GoFieldName(), t, p.Name)
			}
		}
		if part == "" {
			continue
		}
		parts = append(parts, part)
	}
	for _, p := range o.PathParams {
		var part string
		switch {
		case util.IsGoStrType(p.Schema.GoType):
			part = fmt.Sprintf(`t.%s = c.Param("%s")`, p.GoFieldName(), p.Name)
		case util.IsGoNumType(p.Schema.GoType):
			part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.Param("%s")))`, p.GoFieldName(), p.Schema.GoType, p.Name)
		default:
			var t = p.Schema.RefType
			if t == "" {
				t = p.Schema.GoType
			}
			switch {
			case t == "tex.JsInt64" || util.IsGoNumType(p.Schema.Spec.Format):
				part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.Param("%s")))`, p.GoFieldName(), t, p.Name)
			case util.IsGoStrType(p.Schema.Spec.Type):
				part = fmt.Sprintf(`t.%s = %s(c.Param("%s"))`, p.GoFieldName(), t, p.Name)
			}
		}
		if part == "" {
			continue
		}
		parts = append(parts, part)
	}
	for _, p := range o.HeaderParams {
		var part string
		switch {
		case util.IsGoStrType(p.Schema.GoType):
			part = fmt.Sprintf(`t.%s = c.GetHeader("%s")`, p.GoFieldName(), p.Name)
		case util.IsGoNumType(p.Schema.GoType):
			part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.GetHeader("%s")))`, p.GoFieldName(), p.Schema.GoType, p.Name)
		case len(p.Schema.Spec.Enum) > 0:
			switch {
			case util.IsGoStrType(p.Schema.Spec.Type):
				part = fmt.Sprintf(`t.%s = %s(c.GetHeader("%s"))`, p.GoFieldName(), p.Schema.RefType, p.Name)
			case util.IsGoNumType(p.Schema.Spec.Format):
				part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.GetHeader("%s")))`, p.GoFieldName(), p.Schema.RefType, p.Name)
			}
		}
		if part == "" {
			continue
		}
		parts = append(parts, part)
	}
	for _, p := range o.CookieParams {
		var part = fmt.Sprintf(`%s, err := c.Cookie("%s")
if err != nil {
	return err
}
`, p.Name, p.Name)
		switch {
		case util.IsGoStrType(p.Schema.GoType):
			part += fmt.Sprintf(`t.%s = %s`, p.GoFieldName(), p.Name)
		case util.IsGoNumType(p.Schema.GoType):
			part += fmt.Sprintf(`t.%s = %s(tex.ToInt64(%s))`, p.GoFieldName(), p.Schema.GoType, p.Name)
		case len(p.Schema.Spec.Enum) > 0:
			switch {
			case util.IsGoStrType(p.Schema.Spec.Type):
				part += fmt.Sprintf(`t.%s = %s(%s)`, p.GoFieldName(), p.Schema.RefType, p.Name)
			case util.IsGoNumType(p.Schema.Spec.Format):
				part += fmt.Sprintf(`t.%s = %s(tex.ToInt64(%s))`, p.GoFieldName(), p.Schema.RefType, p.Name)
			}
		}
		parts = append(parts, part)
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "\n")
}

func (o Operation) ReqParamNormalize() string {
	return ""
}

func (o Operation) ReqBodyName() string {
	if o.BodyType == nil {
		return ""
	}
	return o.BodyType.Name()
}

func (o Operation) ReqBodyNameWithParam() string {
	if o.Method != http.MethodPost {
		return ""
	}
	var r = o.ReqBodyName()
	if r == "" {
		return ""
	}
	return fmt.Sprintf(", body *%s", r)
}

func (o Operation) ReqBodyBind() string {
	if o.Body == nil {
		return ""
	}
	switch o.Body.ContentType {
	case "application/json":
		return "return c.ShouldBindJSON(t)"
	case "multipart/form-data", "application/x-www-form-urlencoded":
		var parts []string
		for _, p := range o.Body.Schema.Properties {
			var part string
			switch {
			case p.Schema.GoType == "*FileHeader":
				part = fmt.Sprintf(`%s, err := FormGinFile(c, "%s")
if err != nil {
	return err
}
t.%s = %s
`, p.Name, p.Name, p.GoFieldName(), p.Name)
			case p.Schema.GoType == "bool":
				part = fmt.Sprintf(`t.%s = tex.ToBool(c.PostForm("%s"))`, p.GoFieldName(), p.Name)
			case util.IsGoStrType(p.Schema.GoType):
				part = fmt.Sprintf(`t.%s = c.PostForm("%s")`, p.GoFieldName(), p.Name)
			case util.IsGoNumType(p.Schema.GoType):
				part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.PostForm("%s")))`, p.GoFieldName(), p.Schema.GoType, p.Name)
			default:
				var t = p.Schema.RefType
				if t == "" {
					t = p.Schema.GoType
				}
				switch {
				case t == "tex.JsInt64" || util.IsGoNumType(p.Schema.Spec.Format):
					part = fmt.Sprintf(`t.%s = %s(tex.ToInt64(c.PostForm("%s")))`, p.GoFieldName(), t, p.Name)
				case util.IsGoStrType(p.Schema.Spec.Type):
					part = fmt.Sprintf(`t.%s = %s(c.PostForm("%s"))`, p.GoFieldName(), t, p.Name)
				}
			}
			if part == "" {
				continue
			}
			parts = append(parts, part)
		}
		return strings.Join(append(parts, "return nil"), "\n")
	default:
		return ""
	}
}

func (o Operation) ReqBodyNormalize() string {
	return ""
}

func (o Operation) ResponseName() string {
	if o.ServeFile {
		return "FileData"
	}
	if o.ServeHTML {
		return "render.Render"
	}
	if o.ServeText {
		return "[]byte"
	}
	if !o.CreateResponse() {
		return o.ResponseType.Schema.RefType
	}
	return fmt.Sprintf("%sResponse", o.OperationID)
}

func (o Operation) CreateResponse() bool {
	if o.ResponseType != nil && o.ResponseType.Schema != nil {
		return o.ResponseType.Schema.RefType == ""
	}
	return true
}

func (o Operation) GetResponseName() string {
	var r = o.ResponseName()
	if !o.ServeText && !o.ServeHTML {
		return "*" + r
	}
	return r
}

func (o Operation) GetFuncReturn() string {
	if o.NoneLogic {
		return "error"
	}
	return fmt.Sprintf("(%s, error)", o.GetResponseName())
}

func (o Operation) GetEmptyReturn() string {
	if o.NoneLogic {
		return "nil"
	}
	return "nil, nil"
}

func (o Operation) ResponseTypeName() string {
	var t = o.ResponseType.Schema.TypeDecl()
	if !strings.HasPrefix(t, "*") {
		return t
	}
	return t[1:]
}

func (o Operation) NeedSecurity() bool {
	return len(o.Securities) > 0
}

// ParseOperations 解析Op
func ParseOperations(paths openapi3.Paths) ([]Operation, error) {
	var ops = make([]Operation, 0)
	for _, path := range util.SortedPathByKeys(paths) {
		var item = paths[path]
		for method, op := range item.Operations() {
			if op.OperationID == "" {
				return nil, fmt.Errorf("parseOperations.operationID.is.empty (%s)", path)
			}
			var operationID = util.ToCamelCase(op.OperationID)
			var params, err = ParseParams(op.Parameters, []string{operationID + "Params"})
			if err != nil {
				return nil, fmt.Errorf("parseOperations.parseParameters.err (%s)", err.Error())
			}
			var operation = Operation{
				OperationID: operationID,
				Path:        wrapperGinPath(path),
				Method:      method,
				Comment:     util.Str2GoComment(op.Summary, operationID),
				Spec:        op,
			}
			var comment = op.Summary
			if comment == "" {
				comment = op.Description
			}
			operation.Comment = util.Str2GoComment(comment, operationID)
			if op.Security != nil {
				operation.Securities = ParseSecurity(*op.Security)
			}
			if ext, ok := op.Extensions[extOpNoneLogic]; ok {
				operation.NoneLogic, _ = parseExtBool(ext)
			}
			if !operation.NoneLogic {
				operation.HeaderParams = FilterParamsByIn(params, "header")
				operation.QueryParams = FilterParamsByIn(params, "query")
				operation.CookieParams = FilterParamsByIn(params, "cookie")
				operation.PathParams = FilterParamsByIn(params, "path")
				if method == http.MethodPost {
					operation.Body, operation.BodyType, err = ParseRequestBody(operationID, op.RequestBody)
					if err != nil {
						return nil, fmt.Errorf("parseOperations.parseRequestBody.err (%s)", err.Error())
					}
				}
				operation.ParamType = GenerateOpParamType(&operation)
				if ext, ok := op.Extensions[extOpServeFile]; ok {
					operation.ServeFile, _ = parseExtBool(ext)
				}
				if !operation.ServeFile {
					operation.ResponseType, operation.ContentType, err = ParseResponse(op.Responses,
						[]string{operationID + "Response"},
					)
					if err != nil {
						return nil, fmt.Errorf("parseOperations.parseResponse.err (%s)", err.Error())
					}
					switch operation.ContentType {
					case "/*":
						operation.ServeFile = true
					case "text/html":
						operation.ServeHTML = true
					case "text/plain":
						operation.ServeText = true
					default:
						operation.ServeJSON = true
					}
				}
			}
			ops = append(ops, operation)
		}
	}
	return ops, nil
}

// GenerateOpParamType 生成Op的参数Type
func GenerateOpParamType(op *Operation) *Type {
	var params = op.AllParams()
	if len(params) == 0 {
		return nil
	}
	var typeName = op.OperationID + "Params"
	var schema = &Schema{
		Comment: util.Str2GoComment(op.Spec.Description, typeName),
	}
	for _, param := range params {
		schema.Properties = append(schema.Properties, Property{
			Name:     param.Name,
			Comment:  util.Str2GoComment(param.Spec.Description, param.Name),
			Schema:   param.Schema,
			Required: param.Required,
		})
	}
	schema.GoType = GenStructBySchema(typeName, schema)
	return &Type{JSONName: typeName, Schema: schema}
}

// wrapperGinPath 临时处理
func wrapperGinPath(path string) string {
	if !strings.Contains(path, "{") {
		return path
	}
	path = strings.Replace(path, "{", ":", -1)
	path = strings.Replace(path, "}", "", -1)
	return path
}
