package controller

import (
	"github.com/gin-gonic/gin"
	"x-ui/web/session"
)

type BaseController struct {
}

func NewBaseController(g *gin.RouterGroup) *BaseController {
	return &BaseController{}
}

func (a *BaseController) before(c *gin.Context) {
	if !session.IsLogin(c) {
		pureJsonMsg(c, false, "登录时效已过，请重新登录")
		c.Abort()
	} else {
		c.Next()
	}
}
