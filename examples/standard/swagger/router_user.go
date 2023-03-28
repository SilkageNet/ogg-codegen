// Code generated by ogg-codegen. DO NOT EDIT.
package swagger

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/pinealctx/neptune/tex"
	"go.uber.org/zap"
	"net/http"
)

var (
	_ = errors.New
	_ render.Render
	_ tex.JsInt64
)

// 用户类接口

func (w *Wrapper) SetupUserRoutine() {
	w.engine.GET("/user/fetch", w.before("FetchUser"), w.authVerifyFunc, w.FetchUser, w.after("FetchUser"))
	w.engine.POST("/user/set_passwd", w.before("SetPasswd"), w.authVerifyFunc, w.SetPasswd, w.after("SetPasswd"))
	w.engine.POST("/user/update", w.before("UpdateUser"), w.authVerifyFunc, w.UpdateUser, w.after("UpdateUser"))
}

func (w *Wrapper) FetchUser(c *gin.Context) {
	var err error
	if err = w.server.FetchUser(c); err != nil {
		w.logError("FetchUser.err", zap.Error(err))
		w.errorWrapperFunc(c, err)
	}
}

func (m *MockServer) FetchUser(c *gin.Context) error {
	return nil
}

type SetPasswdParams struct {
	Lang    LangCode `json:"lang,omitempty"`
	DevType DevType  `json:"dev_type,omitempty"`
}

func (t *SetPasswdParams) Bind(c *gin.Context) error {
	t.Lang = LangCode(c.Query("lang"))
	t.DevType = DevType(c.Query("dev_type"))
	return nil
}

func (t *SetPasswdParams) Normalize() {

}

func (t *SetPasswdParams) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}

	return nil
}

type SetPasswdBody struct {
	// 密码
	Passwd string `json:"passwd"`
	// 短信验证码
	Code string `json:"code"`
	// 短信验证码HASH
	Hash string `json:"hash"`
}

func (t *SetPasswdBody) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(t)
}

func (t *SetPasswdBody) Normalize() {

}

func (t *SetPasswdBody) Valid() error {
	if t == nil {
		return errors.New("body is nil")
	}
	if t.Passwd == "" {
		return errors.New("passwd.is.empty")
	}
	if t.Code == "" {
		return errors.New("code.is.empty")
	}
	if t.Hash == "" {
		return errors.New("hash.is.empty")
	}
	return nil
}

func (w *Wrapper) SetPasswd(c *gin.Context) {
	var err error
	var param = &SetPasswdParams{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("SetPasswd.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	var body = &SetPasswdBody{}
	if err = w.bindAndValid(c, body); err != nil {
		w.logError("SetPasswd.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	w.logDebug("SetPasswd", zap.Reflect("param", param), zap.Reflect("body", body))
	var res *Response
	if res, err = w.server.SetPasswd(c, param, body); err != nil {
		w.logError("SetPasswd.err", zap.Reflect("param", param), zap.Reflect("body", body), zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("SetPasswd.rsp", zap.Reflect("param", param), zap.Reflect("body", body), zap.Reflect("res", res))
	c.JSON(http.StatusOK, res)
}

func (m *MockServer) SetPasswd(c *gin.Context, params *SetPasswdParams, body *SetPasswdBody) (*Response, error) {
	return nil, nil
}

type UpdateUserParams struct {
	Lang    LangCode `json:"lang,omitempty"`
	DevType DevType  `json:"dev_type,omitempty"`
}

func (t *UpdateUserParams) Bind(c *gin.Context) error {
	t.Lang = LangCode(c.Query("lang"))
	t.DevType = DevType(c.Query("dev_type"))
	return nil
}

func (t *UpdateUserParams) Normalize() {

}

func (t *UpdateUserParams) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}

	return nil
}

type UpdateUserBody struct {
	// 昵称
	Nickname string `json:"nickname,omitempty"`
	// 邮箱
	Email string `json:"email,omitempty"`
	// 微信
	Wechat string `json:"wechat,omitempty"`
	// 头像ID
	AvatarId string `json:"avatar_id,omitempty"`
}

func (t *UpdateUserBody) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(t)
}

func (t *UpdateUserBody) Normalize() {

}

func (t *UpdateUserBody) Valid() error {
	if t == nil {
		return errors.New("body is nil")
	}
	if t.Nickname != "" && len(t.Nickname) < 1 && len(t.Nickname) > 20 {
		return errors.New("nickname.length.error")
	}
	return nil
}

func (w *Wrapper) UpdateUser(c *gin.Context) {
	var err error
	var param = &UpdateUserParams{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("UpdateUser.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	var body = &UpdateUserBody{}
	if err = w.bindAndValid(c, body); err != nil {
		w.logError("UpdateUser.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	w.logDebug("UpdateUser", zap.Reflect("param", param), zap.Reflect("body", body))
	var res *Response
	if res, err = w.server.UpdateUser(c, param, body); err != nil {
		w.logError("UpdateUser.err", zap.Reflect("param", param), zap.Reflect("body", body), zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("UpdateUser.rsp", zap.Reflect("param", param), zap.Reflect("body", body), zap.Reflect("res", res))
	c.JSON(http.StatusOK, res)
}

func (m *MockServer) UpdateUser(c *gin.Context, params *UpdateUserParams, body *UpdateUserBody) (*Response, error) {
	return nil, nil
}
