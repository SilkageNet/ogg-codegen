package codegen

import "github.com/getkin/kin-openapi/openapi3"

func (g *Generator) filterOpsByTag() {
	if len(g.includeTags) == 0 {
		return
	}
	for _, pathItem := range g.swagger.Paths {
		var ops = pathItem.Operations()
		var methods = make([]string, 0, len(ops))
		for method, op := range ops {
			if !checkOpHasTag(op, g.includeTags) {
				methods = append(methods, method)
			}
		}
		for _, method := range methods {
			pathItem.SetOperation(method, nil)
		}
	}
}

func (g *Generator) filterSchemas() {
	if len(g.excludeSchemas) == 0 {
		return
	}
	for _, schema := range g.excludeSchemas {
		delete(g.swagger.Components.Schemas, schema)
	}
}

func checkOpHasTag(op *openapi3.Operation, tags []string) bool {
	if op == nil {
		return false
	}
	for _, hasTag := range op.Tags {
		for _, wantTag := range tags {
			if hasTag == wantTag {
				return true
			}
		}
	}
	return false
}
