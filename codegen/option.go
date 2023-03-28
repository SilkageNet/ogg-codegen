package codegen

import "strings"

var defaultOptions = &options{
	packageName: "swagger",
	outputPath:  ".",
	genSchema:   true,
	genServer:   true,
}

type options struct {
	packageName    string
	outputPath     string
	includeTags    []string
	excludeSchemas []string
	genSchema      bool
	genServer      bool
	imports        map[string]string
}

type OptionsFunc func(*options)

//WithPackageName Package name
func WithPackageName(name string) OptionsFunc {
	return func(o *options) {
		if name != "" {
			o.packageName = name
		}
	}
}

//WithIncludeTags Only include operations that have one of these tags. Ignored when empty.
func WithIncludeTags(tags []string) OptionsFunc {
	return func(o *options) {
		if len(tags) != 0 {
			o.includeTags = tags
		}
	}
}

//WithExcludeSchemas Exclude from generation schemas with given names. Ignored when empty.
func WithExcludeSchemas(schemas []string) OptionsFunc {
	return func(o *options) {
		if len(schemas) != 0 {
			o.excludeSchemas = schemas
		}
	}
}

//WithGenSchema Generator schema
func WithGenSchema(genSchema bool) OptionsFunc {
	return func(o *options) {
		o.genSchema = genSchema
	}
}

//WithGenServer Generator server
func WithGenServer(genServer bool) OptionsFunc {
	return func(o *options) {
		o.genServer = genServer
	}
}

//WithImports With imports
func WithImports(imports []string) OptionsFunc {
	return func(o *options) {
		if len(imports) == 0 {
			return
		}
		o.imports = make(map[string]string)
		for _, i := range imports {
			var parts = strings.Split(i, ":")
			if len(parts) != 2 {
				continue
			}
			o.imports[parts[0]] = parts[1]
		}
	}
}

//WithOutputPath With output path
func WithOutputPath(path string) OptionsFunc {
	return func(o *options) {
		if path == "" {
			return
		}
		o.outputPath = path
	}
}
