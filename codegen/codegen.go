package codegen

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen/tmpl"
	"github.com/SilkageNet/ogg-codegen/codegen/util"
	"github.com/SilkageNet/ogg-codegen/openapi"
	"github.com/getkin/kin-openapi/openapi3"
	"os"
	"path"
	"text/template"
)

var (
	defaultImports = map[string]string{
		"errors":                           "_ = errors.New",
		"github.com/pinealctx/neptune/tex": "_ tex.JsInt64",
		"github.com/gin-gonic/gin/render":  "_ render.Render",
	}
	globalSchemas openapi3.Schemas
	KeepFieldSort bool
	T             *openapi.T
)

func Generate(swagger *openapi3.T, fns ...OptionsFunc) error {
	var generator, err = New(swagger, fns...)
	if err != nil {
		return err
	}
	return generator.Exec()
}

type Generator struct {
	options
	swagger  *openapi3.T
	template *template.Template
}

func New(swagger *openapi3.T, fns ...OptionsFunc) (*Generator, error) {
	var o = defaultOptions
	for _, fn := range fns {
		fn(o)
	}
	var imports = defaultImports
	for k, v := range o.imports {
		imports[k] = v
	}
	tmpl.RegisterOptionsFunc(&tmpl.Options{
		PackageName: o.packageName,
		Imports:     imports,
	})
	var t, err = tmpl.Parse()
	if err != nil {
		return nil, err
	}
	var generator = &Generator{
		options:  *o,
		swagger:  swagger,
		template: t,
	}
	return generator, nil
}

func (g *Generator) Exec() error {
	g.filterSchemas()
	g.filterOpsByTag()

	if err := g.createPackage(); err != nil {
		return err
	}

	if g.genSchema {
		var types, err = CreateTypesBySchemas(g.swagger.Components.Schemas)
		if err != nil {
			return err
		}
		if err = g.genSchemaFiles(types); err != nil {
			return err
		}
	}

	if g.genServer {
		globalSchemas = g.swagger.Components.Schemas
		var ops, err = ParseOperations(g.swagger.Paths)
		if err != nil {
			return err
		}
		if err = g.genServerFiles(ops); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) createPackage() error {
	var filepath = path.Join(g.outputPath, g.packageName)
	var _, err = os.Stat(filepath)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(filepath, os.ModePerm)
}

func (g *Generator) genSchemaFiles(types []Type) error {
	var err error
	var schemas, enums = GroupTypesBySchemaAndEnum(types)
	if err = g.genStructFile(schemas); err != nil {
		return fmt.Errorf("genSchema.genStructFile.err (%s)", err.Error())
	}
	if err = g.genEnumFile(enums); err != nil {
		return fmt.Errorf("genSchema.genEnumFile.err (%s)", err.Error())
	}
	return nil
}

func (g *Generator) genStructFile(schemas []Type) error {
	var buffer, err = util.ExecuteTemplate(g.template, "schema.tmpl", schemas)
	if err != nil {
		return fmt.Errorf("genStructFile.executeTemplate.err (%s)", err.Error())
	}
	return util.FormatAndSaveFile(g.makeOutputFilePath("schema.go"), []byte(buffer))
}

func (g *Generator) genEnumFile(enums []Type) error {
	var buffer, err = util.ExecuteTemplate(g.template, "enum.tmpl", enums)
	if err != nil {
		return fmt.Errorf("genEnumFile.executeTemplate.err (%s)", err.Error())
	}
	return util.FormatAndSaveFile(g.makeOutputFilePath("enum.go"), []byte(buffer))
}

func (g *Generator) genServerFiles(ops []Operation) error {
	var err = g.genWrapper()
	if err != nil {
		return fmt.Errorf("genServerFiles.genWrapper.err (%s)", err.Error())
	}
	var tags = GroupOperationByTag(g.swagger, ops)
	if err = g.genServerFile(tags); err != nil {
		return fmt.Errorf("genServerFiles.genServerFile.err (%s)", err.Error())
	}
	if err = g.genRouterFiles(tags); err != nil {
		return fmt.Errorf("genServerFiles.genRouterFiles.err (%s)", err.Error())
	}
	return nil
}

func (g *Generator) genWrapper() error {
	var buffer, err = util.ExecuteTemplate(g.template, "wrapper.tmpl", nil)
	if err != nil {
		return fmt.Errorf("genWrapper.executeTemplate.err (%s)", err.Error())
	}
	return util.FormatAndSaveFile(g.makeOutputFilePath("wrapper.go"), []byte(buffer))
}

func (g *Generator) genServerFile(tags []Tag) error {
	var buffer, err = util.ExecuteTemplate(g.template, "server.tmpl", tags)
	if err != nil {
		return fmt.Errorf("genServerFile.executeTemplate.err (%s)", err.Error())
	}
	return util.FormatAndSaveFile(g.makeOutputFilePath("server.go"), []byte(buffer))
}

func (g *Generator) genRouterFiles(tags []Tag) error {
	var err error
	for _, tag := range tags {
		if err = g.genRouterFile(tag); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) genRouterFile(tag Tag) error {
	var buffer, err = util.ExecuteTemplate(g.template, "router.tmpl", tag)
	if err != nil {
		return fmt.Errorf("genRouterFile.executeTemplate.err (%s)", err.Error())
	}
	return util.FormatAndSaveFile(g.makeOutputFilePath(fmt.Sprintf("router_%s.go", tag.Name)), []byte(buffer))
}

func (g *Generator) makeOutputFilePath(filename string) string {
	return path.Join(g.outputPath, g.packageName, filename)
}
