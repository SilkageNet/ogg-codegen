package codegen

import (
	"bytes"
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"strings"
)

// Property 属性定义
type Property struct {
	Name     string  // 属性名称
	Comment  string  // 注释
	Schema   *Schema // Schema
	Required bool    // 必填
	Nullable bool    // 可空的
}

func (p Property) GoFieldName() string {
	return util.ToCamelCase(p.Name)
}

func (p Property) GoTypeDef() string {
	return p.Schema.TypeDecl()
}

func (p Property) IsEqual(n Property) bool {
	return p.Name == n.Name && !p.IsSchemaEqual(n)
}

func (p Property) IsSchemaEqual(n Property) bool {
	return p.Name == n.Name && p.Schema.TypeDecl() == n.Schema.TypeDecl() && p.Required == n.Required
}

func (p Property) GenZeroValueCheck() string {
	var (
		fName = p.GoFieldName()
		name  = p.Name
	)
	switch {
	case util.IsGoNumType(p.Schema.GoType) || p.Schema.GoType == "tex.JsInt64":
		return fmt.Sprintf(`if t.%s == 0 { return errors.New("%s.is.zero") }`, fName, name)
	case util.IsGoStrType(p.Schema.GoType):
		return fmt.Sprintf(`if t.%s == "" { return errors.New("%s.is.empty") }`, fName, name)
	case p.Schema.IsEnum():
		return fmt.Sprintf(`if err := t.%s.Valid(); err != nil { return err }`, fName)
	case strings.HasPrefix(p.Schema.TypeDecl(), "*"):
		return fmt.Sprintf(`if t.%s == nil { return errors.New("%s.is.nil") }`, fName, name)
	default:
		return ""
	}
}

func (p Property) GenSchemaCheck() string {
	var (
		fName = p.GoFieldName()
		name  = p.Name
	)
	switch {
	case util.IsGoNumType(p.Schema.GoType) || p.Schema.GoType == "tex.JsInt64":
		var parts = make([]string, 0, 2)
		if p.Schema.Spec.Min != nil {
			parts = append(parts, fmt.Sprintf("t.%s < %d", fName, int64(*p.Schema.Spec.Min)))
		}
		if p.Schema.Spec.Max != nil {
			parts = append(parts, fmt.Sprintf("t.%s > %d", fName, int64(*p.Schema.Spec.Max)))
		}
		if len(parts) == 0 {
			return ""
		}
		return fmt.Sprintf(`if t.%s != 0 && %s { return errors.New("%s.range.error") }`,
			fName, strings.Join(parts, " && "), name)
	case util.IsGoStrType(p.Schema.GoType):
		var parts = make([]string, 0, 2)
		if p.Schema.Spec.MinLength != 0 {
			parts = append(parts, fmt.Sprintf("len(t.%s) < %d", fName, p.Schema.Spec.MinLength))
		}
		if p.Schema.Spec.MaxLength != nil {
			parts = append(parts, fmt.Sprintf("len(t.%s) > %d", fName, *p.Schema.Spec.MaxLength))
		}
		if len(parts) == 0 {
			return ""
		}
		return fmt.Sprintf(`if t.%s != "" && %s { return errors.New("%s.length.error") }`,
			fName, strings.Join(parts, " && "), name)
	case p.Schema.GoType == "*FileHeader":
		var temp = ""
		if len(p.Schema.ExtList) != 0 {
			temp += fmt.Sprintf(`if !t.%s.ValidExts("%s") { return errors.New("%s.ext.invalid") }
`,
				fName, strings.Join(p.Schema.ExtList, `","`), name)
		}
		if p.Schema.Spec.MaxLength != nil && *p.Schema.Spec.MaxLength != 0 {
			temp += fmt.Sprintf(`if !t.%s.ValidSize(%d) { return errors.New("%s.size.too.large") }
`,
				fName, *p.Schema.Spec.MaxLength, name)
		}
		return temp
	case strings.HasPrefix(p.Schema.TypeDecl(), "*"):
		return fmt.Sprintf(`if t.%s != nil && t.%s.Valid() != nil { return errors.New("%s.is.nil") }`,
			fName, fName, name)
	default:
		return ""
	}
}

func (p Property) GenValidCode() string {
	var buf = &bytes.Buffer{}
	var code string
	if p.Required {
		code = p.GenZeroValueCheck()
		if code != "" {
			buf.WriteString(code)
		}
	}
	code = p.GenSchemaCheck()
	if code != "" {
		buf.WriteRune('\n')
		buf.WriteString(code)
	}
	return buf.String()
}

// GenFieldsByProps 根据属性生成结构体字段
func GenFieldsByProps(props []Property) []string {
	var err error
	var fields = make([]string, len(props))
	for i, p := range props {
		var buffer = &bytes.Buffer{}
		if p.Comment != "" {
			buffer.WriteString(p.Comment)
			buffer.WriteRune('\n')
		}
		buffer.WriteString(p.GoFieldName())
		buffer.WriteRune(' ')
		buffer.WriteString(p.GoTypeDef())
		buffer.WriteRune(' ')
		buffer.WriteRune('`')
		// 支持 x-omitempty
		var omitEmpty = true
		if extension, ok := p.Schema.Spec.Extensions[extPropOmitEmpty]; ok {
			omitEmpty, err = parseExtBool(extension)
			if err != nil {
				omitEmpty = true
			}
		}
		var fieldTags = map[string]string{"json": p.Name}
		if !p.Required || !omitEmpty {
			fieldTags["json"] += ",omitempty"
		}
		// 支持 x-go-tags
		if extension, ok := p.Schema.Spec.Extensions[extPropGoTags]; ok {
			var tags map[string]string
			if tags, err = extGoTags(extension); err == nil {
				var keys = util.SortedStringKeys(tags)
				for _, k := range keys {
					fieldTags[k] = tags[k]
				}
			}
		}
		var keys = util.SortedStringKeys(fieldTags)
		var l = len(keys)
		for j, k := range keys {
			buffer.WriteString(k)
			buffer.WriteString(`:"`)
			buffer.WriteString(fieldTags[k])
			buffer.WriteRune('"')
			if j != l-1 {
				buffer.WriteRune(' ')
			}
		}
		buffer.WriteRune('`')
		fields[i] = buffer.String()
	}
	return fields
}
