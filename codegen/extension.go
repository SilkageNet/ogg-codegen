package codegen

import (
	"encoding/json"
	"fmt"
)

const (
	extPropGoType    = "x-go-type"
	extPropGoEnum    = "x-go-enum"
	extPropGoTags    = "x-go-tags"
	extPropGoFileExt = "x-go-file-ext"
	extPropOmitEmpty = "x-omitempty"
	extOpServeFile   = "x-serve-file"
	extOpNoneLogic   = "x-none-logic"
)

func extGoType(value interface{}) (string, error) {
	var raw, ok = value.(json.RawMessage)
	if !ok {
		return "", fmt.Errorf("extGoType.convert.err (%v)", value)
	}
	var name string
	var err = json.Unmarshal(raw, &name)
	if err != nil {
		return "", fmt.Errorf("extGoType.unmarshal.err (%s)", err.Error())
	}
	return name, nil
}

func extGoEnums(value interface{}) (map[string]interface{}, error) {
	var raw, ok = value.(json.RawMessage)
	if !ok {
		return nil, fmt.Errorf("extGoEnums.convert.err (%v)", value)
	}
	var tags map[string]interface{}
	if err := json.Unmarshal(raw, &tags); err != nil {
		return nil, fmt.Errorf("extGoEnums.unmarshal.err (%s)", err.Error())
	}
	return tags, nil
}

func extGoTags(value interface{}) (map[string]string, error) {
	var raw, ok = value.(json.RawMessage)
	if !ok {
		return nil, fmt.Errorf("extGoTags.convert.err (%v)", value)
	}
	var tags map[string]string
	if err := json.Unmarshal(raw, &tags); err != nil {
		return nil, fmt.Errorf("extGoTags.unmarshal.err (%s)", err.Error())
	}
	return tags, nil
}

func parseExtBool(value interface{}) (bool, error) {
	raw, ok := value.(json.RawMessage)
	if !ok {
		return false, fmt.Errorf("parseExtBool.convert.err (%v)", value)
	}
	var omitEmpty bool
	if err := json.Unmarshal(raw, &omitEmpty); err != nil {
		return false, fmt.Errorf("parseExtBool.unmarshal.err (%s)", err.Error())
	}
	return omitEmpty, nil
}

func parseExtSlice(value interface{}) ([]string, error) {
	var raw, ok = value.(json.RawMessage)
	if !ok {
		return nil, fmt.Errorf("parseExtSlice.convert.err (%v)", value)
	}
	var tags []string
	if err := json.Unmarshal(raw, &tags); err != nil {
		return nil, fmt.Errorf("parseExtSlice.unmarshal.err (%s)", err.Error())
	}
	return tags, nil
}
