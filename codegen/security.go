package codegen

import "github.com/getkin/kin-openapi/openapi3"

type Security struct {
	Name   string
	Scopes []string
}

//ParseSecurity 解析Security
func ParseSecurity(security openapi3.SecurityRequirements) []Security {
	var outDefs = make([]Security, 0)
	for _, sr := range security {
		for k, v := range sr {
			outDefs = append(outDefs, Security{Name: k, Scopes: v})
		}
	}
	return outDefs
}
