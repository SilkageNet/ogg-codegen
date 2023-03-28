package tmpl

const wrapperTmpl = `// Code generated by ogg-codegen. DO NOT EDIT.
package {{ opts.PackageName }}

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

const (
	ContextOpID = "OperationId"
)

var (
	defaultOption = &option{
		authVerifyFunc:   func(*gin.Context) {},
		errorWrapperFunc: func(c *gin.Context, err error) {
			c.JSON(http.StatusOK, gin.H{"code": 4000, "errMsg": err.Error()})
		},
		reqErrorWrapperFunc: func(c *gin.Context, err error) {
			c.JSON(http.StatusOK, gin.H{"code": 4000, "errMsg": err.Error()})
		},
	}
)

type ErrorWrapperFunc func(*gin.Context, error)

type option struct {
	authVerifyFunc      gin.HandlerFunc
	errorWrapperFunc    ErrorWrapperFunc
	reqErrorWrapperFunc ErrorWrapperFunc
	beforeFunc          gin.HandlerFunc
	afterFunc           gin.HandlerFunc
	logging             bool
}

type OptionFunc func(o *option)

func WithAuthVerifyFunc(authVerifyFunc gin.HandlerFunc) OptionFunc {
	return func(o *option) {
		if authVerifyFunc != nil {
			o.authVerifyFunc = authVerifyFunc
		}
	}
}

func WithErrorWrapperFunc(errorWrapperFunc ErrorWrapperFunc) OptionFunc {
	return func(o *option) {
		if errorWrapperFunc != nil {
			o.errorWrapperFunc = errorWrapperFunc
		}
	}
}

func WithReqErrorWrapperFunc(errorWrapperFunc ErrorWrapperFunc) OptionFunc {
	return func(o *option) {
		if errorWrapperFunc != nil {
			o.reqErrorWrapperFunc = errorWrapperFunc
		}
	}
}

func WithBeforeFunc(beforeFunc gin.HandlerFunc) OptionFunc {
	return func(o *option) {
		if beforeFunc != nil {
			o.beforeFunc = beforeFunc
		}
	}
}

func WithAfterFunc(afterFunc gin.HandlerFunc) OptionFunc {
	return func(o *option) {
		if afterFunc != nil {
			o.afterFunc = afterFunc
		}
	}
}

func WithLog() OptionFunc {
	return func(o *option) {
		o.logging = true
	}
}

type Param interface {
	Bind(c *gin.Context) error
	Valid() error
	Normalize()
}

type Wrapper struct {
	option
	engine *gin.Engine
	server Server
}

func (w *Wrapper) logError(msg string, fields ...zap.Field) {
	if !w.logging {
		return
	}
	ulog.Error(msg, fields...)
}

func (w *Wrapper) logDebug(msg string, fields ...zap.Field) {
	if !w.logging {
		return
	}
	ulog.Debug(msg, fields...)
}

func (w *Wrapper) bindAndValid(c *gin.Context, request Param) error {
	var err error
	if err = request.Bind(c); err != nil {
		return err
	}
	request.Normalize()
	if err = request.Valid(); err != nil {
		return err
	}
	return nil
}


func (w *Wrapper) before(opID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ContextOpID, opID)
		if w.beforeFunc != nil {
			w.beforeFunc(c)
		}
	}
}

func (w *Wrapper) after(opID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if w.afterFunc != nil {
			w.afterFunc(c)
		}
	}
}

type FileHeader struct {
	header *multipart.FileHeader
	ext    string
	buff   []byte
}

func FormGinFile(c *gin.Context, name string) (*FileHeader, error) {
	var header, err = c.FormFile(name)
	if err != nil {
		return nil, err
	}
	return &FileHeader{header: header}, nil
}

func (f *FileHeader) ValidExts(exts ...string) bool {
	if len(exts) == 0 {
		return true
	}
	if f.ext == "" {
		f.ext = filepath.Ext(f.header.Filename)
	}
	for _, e := range exts {
		if e == f.ext {
			return true
		}
	}
	return false
}

func (f *FileHeader) ValidSize(size int64) bool {
	return size == 0 || size <= f.header.Size
}

func (f *FileHeader) Ext() string {
	if f.ext == "" {
		f.ext = filepath.Ext(f.header.Filename)
	}
	return f.ext
}

func (f *FileHeader) Size() int64 {
	return f.header.Size
}

func (f *FileHeader) Filemame() string {
	return f.header.Filename
}

func (f *FileHeader) Content() ([]byte, error) {
	if f.buff == nil {
		var file, err = f.header.Open()
		if err != nil {
			return nil, err
		}
		defer func() { _ = file.Close() }()
		f.buff, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}
	return f.buff, nil
}

type FileData struct {
	Buff []byte
	Name string
}

func (w *Wrapper) serveFile(c *gin.Context, file *FileData) {
	var ext = filepath.Ext(file.Name)
	var reader = bytes.NewReader(file.Buff)
	var headers = map[string]string{
		"Content-Disposition": fmt.Sprintf({{ tmplSymbol }}attachment; filename="%s"{{ tmplSymbol }}, file.Name),
	}
	c.DataFromReader(http.StatusOK, int64(len(file.Buff)), getContentType(ext), reader, headers)
}

var contentTypeHash = map[string]string{
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/msword",
	".png":  "image/png",
	".jpg":  "mage/jpeg",
}

func getContentType(ext string) string {
	var c, ok = contentTypeHash[ext]
	if ok {
		return c
	}
	return "application/octet-stream"
}

`
