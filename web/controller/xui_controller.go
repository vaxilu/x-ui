package controller

import (
	"github.com/gin-gonic/gin"
)

type XUIController struct {
	BaseController
}

func NewXUIController(g *gin.RouterGroup) *XUIController {
	a := &XUIController{}
	a.initRouter(g)
	return a
}

func (a *XUIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xui")

	g.GET("/", a.index)
	g.GET("/accounts", a.index)
	g.GET("/setting", a.setting)
}

func (a *XUIController) index(c *gin.Context) {
	html(c, "index.html", "系统状态", nil)
}

func (a *XUIController) setting(c *gin.Context) {

}
