package main

import (
	"fmt"
	"github.com/SilkageNet/ogg-codegen/codegen"
	"github.com/SilkageNet/ogg-codegen/openapi"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var app = cli.NewApp()
	app.Name = "ogg-codegen"
	app.Description = "OpenAPI gin golang codegen toolkit."
	app.Version = "1.1.5"
	app.Flags = []cli.Flag{
		&cli.PathFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Specify a path to a OpenAPI 3.0 spec file",
			Value:   "swagger.yaml",
		},
		&cli.StringFlag{
			Name:  "cookie",
			Usage: "When the file path is URL and authentication is required, you can specify cookie through this parameter",
		},
		&cli.PathFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Usage:   "Specify a file output path",
			Value:   ".",
		},
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"p"},
			Usage:   "The package name for generated code",
			Value:   "swagger",
		},
		&cli.StringSliceFlag{
			Name:    "tag",
			Aliases: []string{"t"},
			Usage:   "Specify the tag list to generate",
		},
		&cli.StringSliceFlag{
			Name:    "exSchema",
			Aliases: []string{"es"},
			Usage:   "Specifies a list of schemas that do not need to be generated",
		},
		&cli.BoolFlag{
			Name:    "merge",
			Aliases: []string{"m"},
			Usage:   "Merge files only",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:  "schema",
			Usage: "Specifies whether the schema code needs to be generated",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "server",
			Usage: "Specifies whether the server code needs to be generated",
			Value: true,
		},
		&cli.StringSliceFlag{
			Name:  "import",
			Usage: "Specify the import list to be added additionally",
		},
		&cli.BoolFlag{
			Name:    "keepFieldSort",
			Aliases: []string{"ks"},
			Usage:   "Specifies whether to not sort methods and fields",
			Value:   true,
		},
	}
	app.Action = action
	var err = app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func action(c *cli.Context) error {
	var swg, t, err = loadSwagger(c.Path("file"), c.String("cookie"))
	if err != nil {
		errExit("error loading swagger spec: %s", err)
	}
	if c.Bool("merge") {
		err = mergeSwagger(swg, c.Path("out"))
		if err != nil {
			errExit("error merge swagger spec: %s", err)
		}
		return nil
	}
	codegen.KeepFieldSort = c.Bool("keepFieldSort")
	codegen.T = t
	err = codegen.Generate(swg,
		codegen.WithPackageName(c.String("package")),
		codegen.WithOutputPath(c.Path("out")),
		codegen.WithIncludeTags(c.StringSlice("tag")),
		codegen.WithExcludeSchemas(c.StringSlice("exSchema")),
		codegen.WithGenSchema(c.Bool("schema")),
		codegen.WithGenServer(c.Bool("server")),
		codegen.WithImports(c.StringSlice("import")),
	)
	if err != nil {
		errExit("error generate: %s", err)
	}
	return nil
}

func loadSwagger(filePath, cookies string) (*openapi3.T, *openapi.T, error) {
	var loader = openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	var u, err = url.Parse(filePath)
	var buff []byte
	if err != nil || u.Scheme == "" || u.Host == "" {
		buff, err = ioutil.ReadFile(filePath)
	} else {
		buff, err = loadDataFromURI(filePath, cookies)
	}
	if err != nil {
		return nil, nil, err
	}
	t, err := loader.LoadFromData(buff)
	if err != nil {
		return nil, nil, err
	}
	tt, err := openapi.New(buff)
	if err != nil {
		return nil, nil, err
	}
	return t, tt, nil
}

func loadDataFromURI(url, cookies string) ([]byte, error) {
	var req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	var cookieList = parseCookie(cookies)
	for _, ck := range cookieList {
		req.AddCookie(ck)
	}
	var rsp *http.Response
	rsp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rsp.Body.Close() }()
	return ioutil.ReadAll(rsp.Body)
}

func parseCookie(cookies string) []*http.Cookie {
	if cookies == "" {
		return nil
	}
	var parts = strings.Split(cookies, ";")
	var cookieList = make([]*http.Cookie, 0, len(parts))
	for _, c := range parts {
		var subParts = strings.Split(c, "=")
		if len(subParts) != 2 {
			continue
		}
		cookieList = append(cookieList, &http.Cookie{
			Name:  strings.TrimSpace(subParts[0]),
			Value: strings.TrimSpace(subParts[1]),
			Path:  "/",
		})
	}
	return cookieList
}

func mergeSwagger(t *openapi3.T, outFile string) error {
	if outFile == "" {
		outFile = "./swagger.yaml"
	}
	var err error
	var data []byte
	var ext = filepath.Ext(outFile)
	if ext == ".yaml" {
		data, err = yaml.Marshal(t)
	} else {
		data, err = t.MarshalJSON()
	}
	if err != nil {
		return fmt.Errorf("marshal JSON/YAML error: %v", err)
	}
	err = ioutil.WriteFile(outFile, data, 0644)
	if err != nil {
		return fmt.Errorf("write out file error: %v", err)
	}
	return nil
}
