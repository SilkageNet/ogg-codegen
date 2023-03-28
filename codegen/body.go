package codegen

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/getkin/kin-openapi/openapi3"
)

type RequestBody struct {
	Required    bool
	Schema      *Schema
	ContentType string
}

// ParseRequestBody 解析请求Body
func ParseRequestBody(opID string, bodyRef *openapi3.RequestBodyRef) (*RequestBody, *Type, error) {
	if bodyRef == nil || bodyRef.Value == nil || len(bodyRef.Value.Content) == 0 {
		return nil, nil, nil
	}

	var (
		body         = bodyRef.Value
		bodyTypeName = opID + "Body"
		reqBodyType  *Type
		reqBody      *RequestBody
	)

	for ct, content := range body.Content {
		var bodySchema, err = CreateSchema(content.Schema, []string{bodyTypeName})
		if err != nil {
			return nil, nil, fmt.Errorf("parseRequestBody.createSchema.err (%s)", err.Error())
		}
		if util.IsInternalRef(bodyRef.Ref) {
			bodySchema.RefType, err = util.ConvRef2GoType(bodyRef.Ref)
			if err != nil {
				return nil, nil, fmt.Errorf("parseRequestBody.convRef2GoType.err (%s)", err.Error())
			}
		}
		if bodySchema.RefType != "" {
			return nil, nil, fmt.Errorf("parseRequestBody.unsupported.ref (%s)", bodySchema.RefType)
		}
		reqBody = &RequestBody{
			Required:    body.Required,
			Schema:      bodySchema,
			ContentType: ct,
		}
		reqBodyType = &Type{
			JSONName: bodyTypeName,
			Schema:   bodySchema,
		}
		break
	}
	return reqBody, reqBodyType, nil
}
