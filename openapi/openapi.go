package openapi

import (
	"gopkg.in/yaml.v2"
)

type Components struct {
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

type Schemas map[string]*Schema

type Schema struct {
	Ref        string        `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Properties yaml.MapSlice `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type Paths map[string]*PathItem

type PathItem struct {
	Get  *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Post *Operation `json:"post,omitempty" yaml:"post,omitempty"`
}

type Operation struct {
	OperationID string       `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	RequestBody *RequestBody `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   Responses    `json:"responses" yaml:"responses"`
}

type Responses map[string]*Response

type Response struct {
	Ref     string
	Content Content `json:"content,omitempty" yaml:"content,omitempty"`
}

type RequestBody struct {
	Ref     string
	Content Content `json:"content,omitempty" yaml:"content,omitempty"`
}

type Content map[string]*MediaType

type MediaType struct {
	Schema *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type T struct {
	Components Components `json:"components,omitempty" yaml:"components,omitempty"`
	Paths      Paths      `json:"paths" yaml:"paths"`

	sortHash map[string]yaml.MapSlice

	visitedPathItemRefs map[string]struct{}
	visitedRequestBody  map[*RequestBody]struct{}
	visitedResponse     map[*Response]struct{}
	visitedSchema       map[*Schema]struct{}
}

func New(data []byte) (*T, error) {
	var t T
	var err = yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	t.sortHash = make(map[string]yaml.MapSlice)
	for k, v := range t.Components.Schemas {
		if v.Properties != nil {
			t.sortHash[k] = v.Properties
		}
	}
	for _, v := range t.Paths {
		var n = parseOperation(v.Post)
		for kk, vv := range n {
			t.sortHash[kk] = vv
		}
		n = parseOperation(v.Get)
		for kk, vv := range n {
			t.sortHash[kk] = vv
		}
	}
	return &t, nil
}

func (t *T) SortFieldKeys(key string) []string {
	var v, ok = t.sortHash[key]
	if !ok {
		return nil
	}
	var keys = make([]string, len(v))
	for i, vv := range v {
		keys[i] = vv.Key.(string)
	}
	return keys
}

func parseOperation(op *Operation) map[string]yaml.MapSlice {
	if op == nil {
		return nil
	}
	var out = make(map[string]yaml.MapSlice)
	for _, v := range op.Responses {
		if v.Ref != "" {
			continue
		}
		for _, vv := range v.Content {
			if vv.Schema == nil || vv.Schema.Properties == nil {
				continue
			}
			out[op.OperationID+"Response"] = vv.Schema.Properties
		}
	}
	if op.RequestBody == nil || op.RequestBody.Content == nil {
		return out
	}
	for _, vv := range op.RequestBody.Content {
		if vv.Schema == nil || vv.Schema.Properties == nil {
			continue
		}
		out[op.OperationID+"Body"] = vv.Schema.Properties
	}
	return out
}
