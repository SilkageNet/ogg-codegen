package codegen

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/getkin/kin-openapi/openapi3"
	"sort"
	"strings"
)

// Schema 类型格式定义
type Schema struct {
	GoType  string           // Golang类型
	RefType string           // Ref引用类型
	Comment string           // 注释
	Spec    *openapi3.Schema // OpenAPIv3 original Schema

	//结构体
	Properties []Property // 属性列表

	//数组
	ArrayItem *Schema // 数组项格式定义

	//枚举
	EnumValues     EnumValues // 枚举值定义
	EnumValWrapper string     // 枚举值包裹字符

	//拓展名列表
	ExtList []string
}

// IsArray 是否为数组
func (s *Schema) IsArray() bool {
	return s.ArrayItem != nil
}

// IsEnum 是否为枚举
func (s *Schema) IsEnum() bool {
	return s.EnumValues != nil || (s.Spec != nil && len(s.Spec.Enum) > 0)
}

// TypeDecl 渲染的类型名
func (s *Schema) TypeDecl() string {
	var t = s.GoType
	if s.IsRef() {
		t = s.RefType
		if !s.IsEnum() && !s.IsArray() {
			t = "*" + t
		}
	}
	return t
}

// IsRef 是否为引用类型
func (s *Schema) IsRef() bool {
	return s.RefType != ""
}

// GenValidCode 生成校验代码
func (s *Schema) GenValidCode() string {
	if s.IsEnum() {
		var parts = make([]string, len(s.EnumValues))
		var i = 0
		for _, k := range s.EnumValues {
			parts[i] = fmt.Sprintf("t != %s", k.Name)
			i++
		}
		return fmt.Sprintf(`
if %s {
 return errors.New("unsupported.enum.type")
}`, strings.Join(parts, " && "))
	}
	var parts = make([]string, len(s.Properties))
	for i, p := range s.Properties {
		parts[i] = p.GenValidCode()
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

// MergeProp 将属性合入Schema
func (s *Schema) MergeProp(p Property) error {
	for _, e := range s.Properties {
		if e.IsEqual(p) {
			return fmt.Errorf("mergeProp.isRqual.err: %s", e.Name)
		}
	}
	s.Properties = append(s.Properties, p)
	return nil
}

// CreateSchema 根据schemaRef创建Schema
func CreateSchema(schemaRef *openapi3.SchemaRef, path []string) (*Schema, error) {
	var outSchema = &Schema{GoType: "interface{}"}
	if schemaRef == nil {
		return outSchema, nil
	}

	var err error
	var schema = schemaRef.Value
	outSchema.Comment = util.Str2GoComment(schema.Description, util.Path2TypeName(path))
	outSchema.Spec = schema

	// 处理引用类型
	if util.IsInternalRef(schemaRef.Ref) {
		outSchema.GoType, err = util.ConvRef2GoType(schemaRef.Ref)
		if err != nil {
			return nil, fmt.Errorf("createSchema.convRef2GoType.err (%s)", err.Error())
		}
		outSchema.RefType = outSchema.GoType
		return outSchema, nil
	}

	// 不支持AnyOf和OneOf
	if schema.AnyOf != nil || schema.OneOf != nil {
		return outSchema, nil
	}

	// 对AllOf的支持
	if schema.AllOf != nil {
		outSchema, err = MergeSchemas(schema.AllOf, path)
		if err != nil {
			return nil, fmt.Errorf("createSchema.mergeSchemas.err (%s)", err.Error())
		}
		outSchema.Spec = schema
		return outSchema, nil
	}

	// 对拓展类型的支持
	var extension, ok = schema.Extensions[extPropGoType]
	if ok {
		outSchema.GoType, err = extGoType(extension)
		if err != nil {
			return nil, fmt.Errorf("createSchema.extGoType.err (%s)", err.Error())
		}
		return outSchema, nil
	}
	// 支持 x-go-file-ext
	extension, ok = schema.Extensions[extPropGoFileExt]
	if ok {
		var extList []string
		if extList, err = parseExtSlice(extension); err == nil {
			outSchema.ExtList = extList
		}
	}

	var t = schema.Type
	// 对Object的处理
	if t == "" || t == "object" {
		if len(schema.Properties) == 0 {
			if t == "object" {
				outSchema.GoType = "map[string]interface{}"
			}
			return outSchema, nil
		}
		// 对Object的属性处理
		for _, pName := range util.SortedSchemaKeys(schema.Properties) {
			var p = schema.Properties[pName]
			var pPath = append(path, pName)
			var pSchema *Schema
			pSchema, err = CreateSchema(p, pPath)
			if err != nil {
				return nil, fmt.Errorf("createSchema.property.createSchema.err (%s)", err.Error())
			}
			var prop = Property{
				Name:     pName,
				Schema:   pSchema,
				Required: util.StrInArray(pName, schema.Required),
			}
			if p.Value != nil {
				prop.Comment = util.Str2GoComment(p.Value.Description, "")
				prop.Nullable = p.Value.Nullable
			}
			outSchema.Properties = append(outSchema.Properties, prop)
		}
		outSchema.GoType = GenStructBySchema(strings.Join(path, ""), outSchema)
		return outSchema, nil
	}

	// 处理Schema的类型转换
	if err = resolveType(schema, path, outSchema); err != nil {
		return nil, fmt.Errorf("createSchema.resolveType.err (%s)", err.Error())
	}

	// 对枚举的支持
	var enumExtension interface{}
	enumExtension, ok = schema.Extensions[extPropGoEnum]
	if len(schema.Enum) > 0 && ok {
		var extMap map[string]interface{}
		if extMap, err = extGoEnums(enumExtension); err != nil {
			return nil, fmt.Errorf("createSchema.extGoEnums.err (%s)", err.Error())
		}
		outSchema.EnumValues = make([]EnumValue, 0, len(extMap))
		for k, v := range extMap {
			outSchema.EnumValues = append(outSchema.EnumValues, EnumValue{
				Name:  k,
				Value: fmt.Sprintf("%v", v),
			})
		}
		sort.Sort(outSchema.EnumValues)
		if outSchema.GoType == "string" {
			outSchema.EnumValWrapper = `"`
		}
		return outSchema, nil
	}

	return outSchema, nil
}

// MergeSchemas 合并allOf为一个Schema
func MergeSchemas(allOf []*openapi3.SchemaRef, path []string) (*Schema, error) {
	var err error
	var outSchema = &Schema{}
	for _, schemaRef := range allOf {
		var schema *Schema
		schema, err = CreateSchema(schemaRef, path)
		if err != nil {
			return nil, fmt.Errorf("mergeSchemas.createSchema.err (%s)", err.Error())
		}
		if util.IsInternalRef(schemaRef.Ref) {
			schema.RefType, err = util.ConvRef2GoType(schemaRef.Ref)
			if err != nil {
				return nil, fmt.Errorf("mergeSchemas.convRef2GoType.err (%s)", err.Error())
			}
		}
		for _, p := range schema.Properties {
			err = outSchema.MergeProp(p)
			if err != nil {
				return nil, fmt.Errorf("mergeSchemas.mergeProp.err (%s)", err.Error())
			}
		}
	}
	outSchema.GoType, err = GenStructByAllOf(allOf, path)
	if err != nil {
		return nil, fmt.Errorf("mergeSchemas.genStructByAllOf.err (%s)", err.Error())
	}
	return outSchema, nil
}

// GenStructByAllOf 根据allOf生成结构体
func GenStructByAllOf(allOf []*openapi3.SchemaRef, path []string) (string, error) {
	var parts = []string{"struct {"}
	var err error
	for _, schemaRef := range allOf {
		if util.IsInternalRef(schemaRef.Ref) {
			var goType string
			goType, err = util.ConvRef2GoType(schemaRef.Ref)
			if err != nil {
				return "", fmt.Errorf("genStructByAllOf.convRef2GoType.err (%s)", err.Error())
			}
			parts = append(parts, fmt.Sprintf("%s\n", goType))
			continue
		}
		var schema *Schema
		schema, err = CreateSchema(schemaRef, path)
		if err != nil {
			return "", fmt.Errorf("genStructByAllOf.createSchema.err (%s)", err.Error())
		}
		var properties = schema.Properties
		if KeepFieldSort {
			var hash = make(map[string]Property)
			for _, p := range properties {
				hash[p.Name] = p
			}
			var fields = T.SortFieldKeys(strings.Join(path, ""))
			for i, f := range fields {
				properties[i] = hash[f]
			}
		}
		parts = append(parts, GenFieldsByProps(properties)...)
	}
	parts = append(parts, "}")
	return strings.Join(parts, "\n"), nil
}

// GenStructBySchema 根据Schema生成结构体
func GenStructBySchema(name string, schema *Schema) string {
	var properties = schema.Properties
	if KeepFieldSort {
		var hash = make(map[string]Property)
		for _, p := range properties {
			hash[p.Name] = p
		}
		var fields = T.SortFieldKeys(name)
		for i, f := range fields {
			properties[i] = hash[f]
		}
	}
	var parts = []string{"struct {"}
	parts = append(parts, GenFieldsByProps(properties)...)
	parts = append(parts, "}")
	return strings.Join(parts, "\n")
}
