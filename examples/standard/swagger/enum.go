// Code generated by ogg-codegen. DO NOT EDIT.
package swagger

import (
	"errors"
	"github.com/gin-gonic/gin/render"
	"github.com/pinealctx/neptune/tex"
)

var (
	_ = errors.New
	_ render.Render
	_ tex.JsInt64
)

// DevType 设备类型；默认 `web`；所有请求都可以携带该参数，便于请求跟踪与溯源。
// - web：网页
// - mini：小程序
type DevType string

func (t DevType) Valid() error {
	var err error
	if t != DevTypeMini && t != DevTypeWeb {
		return errors.New("unsupported.enum.type")
	}
	return err
}

const (
	DevTypeMini DevType = "mini"
	DevTypeWeb  DevType = "web"
)

// FileType 上传文件类型：
// - 0：头像
// - 1：附件
type FileType uint8

func (t FileType) Valid() error {
	var err error
	if t != FileTypeAvatar && t != FileTypeAnnex {
		return errors.New("unsupported.enum.type")
	}
	return err
}

const (
	FileTypeAvatar FileType = 0
	FileTypeAnnex  FileType = 1
)

// LangCode 语言类型；默认 `zh`；所有请求都可以携带该参数，服务端会根据语言返回相应的数据。
// - zh：中文
// - en：英文
type LangCode string

func (t LangCode) Valid() error {
	var err error
	if t != LangCodeEN && t != LangCodeZH {
		return errors.New("unsupported.enum.type")
	}
	return err
}

const (
	LangCodeEN LangCode = "en"
	LangCodeZH LangCode = "zh"
)

// SendCodeType 获取验证码类型：
// - 0：登录
// - 1：设置密码
type SendCodeType uint8

func (t SendCodeType) Valid() error {
	var err error
	if t != SendCodeTypeLogin && t != SendCodeTypeSetPwd {
		return errors.New("unsupported.enum.type")
	}
	return err
}

const (
	SendCodeTypeLogin  SendCodeType = 0
	SendCodeTypeSetPwd SendCodeType = 1
)

// UserState 用户状态
type UserState uint8

func (t UserState) Valid() error {
	var err error
	if t != UserStateEnabled && t != UserStateDisabled {
		return errors.New("unsupported.enum.type")
	}
	return err
}

const (
	UserStateEnabled  UserState = 0
	UserStateDisabled UserState = 1
)
