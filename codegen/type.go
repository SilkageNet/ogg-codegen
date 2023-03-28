package codegen

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/getkin/kin-openapi/openapi3"
)

// Type 类型定义
type Type struct {
	JSONName string  // JSON名称
	Schema   *Schema // Schema
}

func (t Type) Name() string {
	return util.ToCamelCase(t.JSONName)
}

// CreateTypesBySchemas 根据schemas创建Types
func CreateTypesBySchemas(schemas map[string]*openapi3.SchemaRef) ([]Type, error) {
	var types = make([]Type, 0, len(schemas))
	for _, name := range util.SortedSchemaKeys(schemas) {
		var schemaRef = schemas[name]
		var schema, err = CreateSchema(schemaRef, []string{name})
		if err != nil {
			return nil, fmt.Errorf("createTypesBySchemas.createSchema.err (%s)", err.Error())
		}
		types = append(types, Type{
			JSONName: name,
			Schema:   schema,
		})
	}
	return types, nil
}

// GroupTypesBySchemaAndEnum 将Types按Schema和Enum分组
func GroupTypesBySchemaAndEnum(types []Type) (schemas []Type, enums []Type) {
	var enumHash = make(map[string]Type)
	for _, t := range types {
		if t.Schema.IsEnum() {
			enumHash[t.Name()] = t
			enums = append(enums, t)
		} else {
			schemas = append(schemas, t)
		}
	}

	for i, t := range schemas {
		for ii, p := range t.Schema.Properties {
			var tt, ok = enumHash[p.Schema.RefType]
			if !ok {
				continue
			}
			p.Schema = tt.Schema
			t.Schema.Properties[ii] = p
		}
		schemas[i] = t
	}
	return
}
