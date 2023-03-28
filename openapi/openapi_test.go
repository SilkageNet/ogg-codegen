package openapi

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func Test_Parse(t *testing.T) {
	var data, err = ioutil.ReadFile("../examples/standard/swagger.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	var s T
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s)
}
