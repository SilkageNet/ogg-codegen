package tmpl

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_ParseWrapper(t *testing.T) {
	RegisterOptionsFunc(&Options{PackageName: "swagger"})
	var tmpl, err = Parse()
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	var writer = bufio.NewWriter(&buf)
	if err = tmpl.ExecuteTemplate(writer, "wrapper.tmpl", nil); err != nil {
		t.Error(err)
		return
	}
	if err = writer.Flush(); err != nil {
		t.Error(err)
		return
	}
	t.Log(buf.String())
}
