package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"x-ui/database/model"
	"x-ui/web/service"
	"x-ui/web/session"
)

type XUIController struct {
	BaseController

	inboundService service.InboundService
}

func NewXUIController(g *gin.RouterGroup) *XUIController {
	a := &XUIController{}
	a.initRouter(g)
	return a
}

func (a *XUIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xui")

	g.GET("/", a.index)
	g.GET("/inbounds", a.inbounds)
	g.POST("/inbounds", a.postInbounds)
	g.POST("/inbound/add", a.addInbound)
	g.POST("/inbound/del/:id", a.delInbound)
	g.GET("/setting", a.setting)
}

func (a *XUIController) index(c *gin.Context) {
	html(c, "index.html", "系统状态", nil)
}

func (a *XUIController) inbounds(c *gin.Context) {
	html(c, "inbounds.html", "入站列表", nil)
}

func (a *XUIController) setting(c *gin.Context) {
	html(c, "setting.html", "设置", nil)
}

func (a *XUIController) postInbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	inbounds, err := a.inboundService.GetInbounds(user.Id)
	if err != nil {
		jsonMsg(c, "获取", err)
		return
	}
	jsonObj(c, inbounds, nil)
}

func (a *XUIController) addInbound(c *gin.Context) {
	inbound := &model.Inbound{}
	err := c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "添加", err)
		return
	}
	user := session.GetLoginUser(c)
	inbound.UserId = user.Id
	inbound.Enable = true
	log.Println(inbound)
	err = a.inboundService.AddInbound(inbound)
	jsonMsg(c, "添加", err)
}

func (a *XUIController) delInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "删除", err)
		return
	}
	err = a.inboundService.DelInbound(id)
	jsonMsg(c, "删除", err)
}
