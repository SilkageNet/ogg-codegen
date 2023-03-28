package codegen

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
)

const (
	integerDefault = "int"
	numberDefault  = "float32"
	stringDefault  = "string"
)

var (
	integerHash = map[string]string{
		"int64":  "int64",
		"int32":  "int32",
		"int16":  "int16",
		"int8":   "int8",
		"int":    "int",
		"uint64": "uint64",
		"uint32": "uint32",
		"uint16": "uint16",
		"uint8":  "uint8",
		"uint":   "uint",
	}
	numberHash = map[string]string{
		"double": "float64",
		"float":  "float32",
	}
	stringHash = map[string]string{
		"byte":      "[]byte",
		"date-time": "time.Time",
		"json":      "json.RawMessage",
		"binary":    "*FileHeader",
	}
)

//resolveType 处理Schema的类型为Golang的Type
func resolveType(schema *openapi3.Schema, path []string, outSchema *Schema) error {
	var (
		ok  bool
		err error
		t   = schema.Type
		f   = schema.Format
	)

	switch t {
	case "array":
		var arrayItem *Schema
		arrayItem, err = CreateSchema(schema.Items, path)
		if err != nil {
			return fmt.Errorf("resolveType.array.createSchema.err (%s)", err.Error())
		}
		outSchema.ArrayItem = arrayItem
		outSchema.GoType = "[]" + arrayItem.TypeDecl()
	case "integer":
		outSchema.GoType, ok = integerHash[f]
		if !ok {
			outSchema.GoType = integerDefault
		}
	case "number":
		outSchema.GoType, ok = numberHash[f]
		if !ok {
			outSchema.GoType = numberDefault
		}
	case "boolean":
		outSchema.GoType = "bool"
	case "string":
		outSchema.GoType, ok = stringHash[f]
		if !ok {
			outSchema.GoType = stringDefault
		}
	default:
		return fmt.Errorf("resolveType.type.unsupported.err (%s)", t)
	}
	return nil
}
