package main

import (
	"github.com/SilkageNet/ogg-codegen/examples/standard/swagger"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	var engine = gin.New()
	var srvImpl = &Server{}
	swagger.Register(engine, srvImpl,
		swagger.WithAuthVerifyFunc(srvImpl.AuthVerify),
		swagger.WithErrorWrapperFunc(func(c *gin.Context, err error) {
			c.JSON(http.StatusOK, gin.H{"code": 4000, "errMsg": err.Error()})
		}))
	var srv = &http.Server{Addr: ":8989", Handler: engine}
	var err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

type Server struct {
	swagger.MockServer
}

func (s *Server) PasswdLogin(c *gin.Context, params *swagger.PasswdLoginParams, body *swagger.PasswdLoginBody) (*swagger.PasswdLoginResponse, error) {
	return &swagger.PasswdLoginResponse{
		Response: swagger.Response{
			Code:   2000,
			ErrMsg: "no error",
		},
		Data: &swagger.LoginResult{
			New:   true,
			Token: "this is token",
			User: &swagger.User{
				AccessHash: 12015,
				AreaCode:   "+86",
				Id:         12015,
				Nickname:   "Silkage",
				Phone:      "18888888888",
				State:      swagger.UserStateEnabled,
			},
		},
	}, nil
}

func (s *Server) AuthVerify(c *gin.Context) {

}
