package codegen

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/getkin/kin-openapi/openapi3"
)

//ParseResponse 解析Response
func ParseResponse(responses openapi3.Responses, path []string) (*Type, string, error) {
	for _, res := range responses {
		if res.Value == nil || res.Value.Content == nil {
			continue
		}
		for contentType, item := range res.Value.Content {
			if item.Schema == nil {
				continue
			}
			var resSchema, err = CreateSchema(item.Schema, path)
			if err != nil {
				return nil, "", fmt.Errorf("parseResponse.createSchema.err (%s)", err.Error())
			}
			var response = &Type{
				JSONName: util.Path2TypeName(path),
				Schema:   resSchema,
			}
			if util.IsInternalRef(item.Schema.Ref) {
				response.Schema.RefType, err = util.ConvRef2GoType(item.Schema.Ref)
				if err != nil {
					return nil, "", fmt.Errorf("parseResponse.convRef2GoType.err (%s)", err.Error())
				}
			}
			return response, contentType, nil
		}
	}
	return nil, "", nil
}
