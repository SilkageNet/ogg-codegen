package codegen

import (
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/getkin/kin-openapi/openapi3"
)

type Tag struct {
	Name       string
	Comment    string
	Operations []Operation
	Spec       *openapi3.Tag
}

func (t Tag) GoName() string {
	return util.ToCamelCase(t.Name)
}

//GroupOperationByTag 根据Tag将Op分组
func GroupOperationByTag(swagger *openapi3.T, ops []Operation) []Tag {
	var tags = make([]Tag, len(swagger.Tags))
	var tagHash = make(map[string]int)
	for i, tag := range swagger.Tags {
		tags[i] = Tag{
			Name:    tag.Name,
			Comment: util.Str2GoComment(tag.Description, ""),
			Spec:    tag,
		}
		tagHash[tag.Name] = i
	}
	for _, op := range ops {
		if op.Spec == nil || len(op.Spec.Tags) == 0 {
			continue
		}
		var tag = op.Spec.Tags[0]
		var index, ok = tagHash[tag]
		if !ok {
			continue
		}
		tags[index].Operations = append(tags[index].Operations, op)
	}
	var outTags = make([]Tag, 0, len(swagger.Tags))
	for _, t := range tags {
		if len(t.Operations) == 0 {
			continue
		}
		outTags = append(outTags, t)
	}
	return outTags
}
