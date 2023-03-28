package codegen

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/getkin/kin-openapi/openapi3"
)

type Parameter struct {
	Name     string
	In       string
	Required bool
	Spec     *openapi3.Parameter
	Schema   *Schema
}

func (p Parameter) GoFieldName() string {
	return util.ToCamelCase(p.Name)
}

// ParseParams 解析swagger的params
func ParseParams(params openapi3.Parameters, path []string) ([]Parameter, error) {
	var outParams = make([]Parameter, len(params))
	var err error
	for i, paramRef := range params {
		var param = paramRef.Value
		var p = Parameter{
			Name:     param.Name,
			In:       param.In,
			Required: param.Required,
			Spec:     param,
		}
		p.Schema, err = CreateSchemaByParam(param, append(path, param.Name))
		if err != nil {
			return nil, fmt.Errorf("parseParameters.createSchemaByParam.err (%s)", err.Error())
		}
		var r = p.Schema.RefType
		if r != "" {
			var sc, ok = globalSchemas[r]
			if !ok {
				return nil, fmt.Errorf("parseParameters.schema.ref.not.found (%s)", p.Schema.RefType)
			}
			p.Schema, err = CreateSchema(sc, append(path, param.Name))
			if err != nil {
				return nil, fmt.Errorf("parseParameters.ref.createSchemaByParam.err (%s)", err.Error())
			}
			// 枚举类型依然保留
			if p.Schema.EnumValues.Len() != 0 {
				p.Schema.GoType = r
			}
		}
		outParams[i] = p
	}
	return outParams, nil
}

// CreateSchemaByParam 根据Param创建Schema
func CreateSchemaByParam(param *openapi3.Parameter, path []string) (*Schema, error) {
	if param.Schema == nil {
		return nil, fmt.Errorf("createSchemaByParam.param.invalid (%s)", param.Name)
	}
	return CreateSchema(param.Schema, path)
}

// FilterParamsByIn 根据in类型过滤参数
func FilterParamsByIn(params []Parameter, in string) []Parameter {
	var out []Parameter
	for _, p := range params {
		if p.In == in {
			out = append(out, p)
		}
	}
	return out
}
