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

// 鉴权类接口

func (w *Wrapper) SetupAuthRoutine() {
	w.engine.GET("/auth/:name/:id", w.before("AuthTest"), w.AuthTest, w.after("AuthTest"))
	w.engine.GET("/auth/logout", w.before("Logout"), w.authVerifyFunc, w.Logout, w.after("Logout"))
	w.engine.POST("/auth/passwd_login", w.before("PasswdLogin"), w.PasswdLogin, w.after("PasswdLogin"))
	w.engine.POST("/auth/send_code", w.before("SendCode"), w.SendCode, w.after("SendCode"))
	w.engine.POST("/auth/sms_login", w.before("SMSLogin"), w.SMSLogin, w.after("SMSLogin"))
}

type AuthTestParams struct {
	Lang    LangCode `json:"lang,omitempty"`
	DevType DevType  `json:"dev_type,omitempty"`
	Test    string   `json:"test,omitempty"`
	// name 名称
	Name string `json:"name,omitempty"`
	// id ID
	Id int32 `json:"id,omitempty"`
}

func (t *AuthTestParams) Bind(c *gin.Context) error {
	t.Lang = LangCode(c.Query("lang"))
	t.DevType = DevType(c.Query("dev_type"))
	t.Test = c.Query("test")
	t.Name = c.Param("name")
	t.Id = int32(tex.ToInt64(c.Param("id")))
	return nil
}

func (t *AuthTestParams) Normalize() {

}

func (t *AuthTestParams) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}

	return nil
}

func (w *Wrapper) AuthTest(c *gin.Context) {
	var err error
	var param = &AuthTestParams{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("AuthTest.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	w.logDebug("AuthTest", zap.Reflect("param", param))
	var res *Response
	if res, err = w.server.AuthTest(c, param); err != nil {
		w.logError("AuthTest.err", zap.Reflect("param", param), zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("AuthTest.rsp", zap.Reflect("param", param), zap.Reflect("res", res))
	c.JSON(http.StatusOK, res)
}

func (m *MockServer) AuthTest(c *gin.Context, params *AuthTestParams) (*Response, error) {
	return nil, nil
}

type LogoutParams struct {
	Lang    LangCode `json:"lang,omitempty"`
	DevType DevType  `json:"dev_type,omitempty"`
}

func (t *LogoutParams) Bind(c *gin.Context) error {
	t.Lang = LangCode(c.Query("lang"))
	t.DevType = DevType(c.Query("dev_type"))
	return nil
}

func (t *LogoutParams) Normalize() {

}

func (t *LogoutParams) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}

	return nil
}

func (w *Wrapper) Logout(c *gin.Context) {
	var err error
	var param = &LogoutParams{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("Logout.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	w.logDebug("Logout", zap.Reflect("param", param))
	var res *Response
	if res, err = w.server.Logout(c, param); err != nil {
		w.logError("Logout.err", zap.Reflect("param", param), zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("Logout.rsp", zap.Reflect("param", param), zap.Reflect("res", res))
	c.JSON(http.StatusOK, res)
}

func (m *MockServer) Logout(c *gin.Context, params *LogoutParams) (*Response, error) {
	return nil, nil
}

type PasswdLoginParams struct {
	Lang    LangCode `json:"lang,omitempty"`
	DevType DevType  `json:"dev_type,omitempty"`
}

func (t *PasswdLoginParams) Bind(c *gin.Context) error {
	t.Lang = LangCode(c.Query("lang"))
	t.DevType = DevType(c.Query("dev_type"))
	return nil
}

func (t *PasswdLoginParams) Normalize() {

}

func (t *PasswdLoginParams) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}

	return nil
}

type PasswdLoginBody struct {
	// 地区码
	AreaCode string `json:"area_code"`
	// 手机号码
	Phone string `json:"phone"`
	// 密码
	Passwd string `json:"passwd"`
	// 验证码通行证
	Ticket string `json:"ticket"`
	// 验证码随机串(非微信小程序模式下必传)
	Randstr string `json:"randstr,omitempty"`
}

func (t *PasswdLoginBody) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(t)
}

func (t *PasswdLoginBody) Normalize() {

}

func (t *PasswdLoginBody) Valid() error {
	if t == nil {
		return errors.New("body is nil")
	}
	if t.AreaCode == "" {
		return errors.New("area_code.is.empty")
	}
	if t.Phone == "" {
		return errors.New("phone.is.empty")
	}
	if t.Passwd == "" {
		return errors.New("passwd.is.empty")
	}
	if t.Ticket == "" {
		return errors.New("ticket.is.empty")
	}
	return nil
}

type PasswdLoginResponse struct {
	Response

	// 登录结果
	Data *LoginResult `json:"data,omitempty"`
}

func (w *Wrapper) PasswdLogin(c *gin.Context) {
	var err error
	var param = &PasswdLoginParams{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("PasswdLogin.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	var body = &PasswdLoginBody{}
	if err = w.bindAndValid(c, body); err != nil {
		w.logError("PasswdLogin.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	w.logDebug("PasswdLogin", zap.Reflect("param", param), zap.Reflect("body", body))
	var res *PasswdLoginResponse
	if res, err = w.server.PasswdLogin(c, param, body); err != nil {
		w.logError("PasswdLogin.err", zap.Reflect("param", param), zap.Reflect("body", body), zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("PasswdLogin.rsp", zap.Reflect("param", param), zap.Reflect("body", body), zap.Reflect("res", res))
	c.JSON(http.StatusOK, res)
}

func (m *MockServer) PasswdLogin(c *gin.Context, params *PasswdLoginParams, body *PasswdLoginBody) (*PasswdLoginResponse, error) {
	return nil, nil
}

type SendCodeParams struct {
	Lang    LangCode `json:"lang,omitempty"`
	DevType DevType  `json:"dev_type,omitempty"`
}

func (t *SendCodeParams) Bind(c *gin.Context) error {
	t.Lang = LangCode(c.Query("lang"))
	t.DevType = DevType(c.Query("dev_type"))
	return nil
}

func (t *SendCodeParams) Normalize() {

}

func (t *SendCodeParams) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}

	return nil
}

type SendCodeBody struct {
	// 获取验证码类型：
	// - 0：登录
	// - 1：设置密码
	Type SendCodeType `json:"type"`
	// 地区码
	AreaCode string `json:"area_code"`
	// 手机号码
	Phone string `json:"phone"`
	// 验证码通行证
	Ticket string `json:"ticket"`
	// 验证码随机串(非微信小程序模式下必传)
	Randstr string `json:"randstr,omitempty"`
}

func (t *SendCodeBody) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(t)
}

func (t *SendCodeBody) Normalize() {

}

func (t *SendCodeBody) Valid() error {
	if t == nil {
		return errors.New("body is nil")
	}
	if err := t.Type.Valid(); err != nil {
		return err
	}
	if t.AreaCode == "" {
		return errors.New("area_code.is.empty")
	}
	if t.Phone == "" {
		return errors.New("phone.is.empty")
	}
	if t.Ticket == "" {
		return errors.New("ticket.is.empty")
	}
	return nil
}

type SendCodeResponse struct {
	Response

	// 验证码唯一HASH标识，在使用验证码的场景需携带。
	Data string `json:"data,omitempty"`
}

func (w *Wrapper) SendCode(c *gin.Context) {
	var err error
	var param = &SendCodeParams{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("SendCode.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	var body = &SendCodeBody{}
	if err = w.bindAndValid(c, body); err != nil {
		w.logError("SendCode.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	w.logDebug("SendCode", zap.Reflect("param", param), zap.Reflect("body", body))
	var res *SendCodeResponse
	if res, err = w.server.SendCode(c, param, body); err != nil {
		w.logError("SendCode.err", zap.Reflect("param", param), zap.Reflect("body", body), zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("SendCode.rsp", zap.Reflect("param", param), zap.Reflect("body", body), zap.Reflect("res", res))
	c.JSON(http.StatusOK, res)
}

func (m *MockServer) SendCode(c *gin.Context, params *SendCodeParams, body *SendCodeBody) (*SendCodeResponse, error) {
	return nil, nil
}

type SMSLoginParams struct {
	Lang    LangCode `json:"lang,omitempty"`
	DevType DevType  `json:"dev_type,omitempty"`
}

func (t *SMSLoginParams) Bind(c *gin.Context) error {
	t.Lang = LangCode(c.Query("lang"))
	t.DevType = DevType(c.Query("dev_type"))
	return nil
}

func (t *SMSLoginParams) Normalize() {

}

func (t *SMSLoginParams) Valid() error {
	if t == nil {
		return errors.New("params is nil")
	}

	return nil
}

type SMSLoginBody struct {
	// 地区码
	AreaCode string `json:"area_code"`
	// 手机号码
	Phone string `json:"phone"`
	// 短信验证码
	Code string `json:"code"`
	// 短信验证码HASHs
	Hash string `json:"hash"`
}

func (t *SMSLoginBody) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(t)
}

func (t *SMSLoginBody) Normalize() {

}

func (t *SMSLoginBody) Valid() error {
	if t == nil {
		return errors.New("body is nil")
	}
	if t.AreaCode == "" {
		return errors.New("area_code.is.empty")
	}
	if t.Phone == "" {
		return errors.New("phone.is.empty")
	}
	if t.Code == "" {
		return errors.New("code.is.empty")
	}
	if t.Hash == "" {
		return errors.New("hash.is.empty")
	}
	return nil
}

type SMSLoginResponse struct {
	Response

	// 登录结果
	Data *LoginResult `json:"data,omitempty"`
}

func (w *Wrapper) SMSLogin(c *gin.Context) {
	var err error
	var param = &SMSLoginParams{}
	if err = w.bindAndValid(c, param); err != nil {
		w.logError("SMSLogin.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	var body = &SMSLoginBody{}
	if err = w.bindAndValid(c, body); err != nil {
		w.logError("SMSLogin.bindAndValid.err", zap.Error(err))
		w.reqErrorWrapperFunc(c, err)
		return
	}
	w.logDebug("SMSLogin", zap.Reflect("param", param), zap.Reflect("body", body))
	var res *SMSLoginResponse
	if res, err = w.server.SMSLogin(c, param, body); err != nil {
		w.logError("SMSLogin.err", zap.Reflect("param", param), zap.Reflect("body", body), zap.Error(err))
		w.errorWrapperFunc(c, err)
		return
	}
	w.logDebug("SMSLogin.rsp", zap.Reflect("param", param), zap.Reflect("body", body), zap.Reflect("res", res))
	c.JSON(http.StatusOK, res)
}

func (m *MockServer) SMSLogin(c *gin.Context, params *SMSLoginParams, body *SMSLoginBody) (*SMSLoginResponse, error) {
	return nil, nil
}