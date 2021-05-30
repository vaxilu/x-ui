package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"log"
	"strconv"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/web/entity"
	"x-ui/web/global"
	"x-ui/web/service"
	"x-ui/web/session"
)

type XUIController struct {
	BaseController

	inboundService service.InboundService
	xrayService    service.XrayService
	settingService service.SettingService

	isNeedXrayRestart atomic.Bool
}

func NewXUIController(g *gin.RouterGroup) *XUIController {
	a := &XUIController{}
	a.initRouter(g)
	a.startTask()
	return a
}

func (a *XUIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xui")
	g.Use(a.checkLogin)

	g.GET("/", a.index)
	g.GET("/inbounds", a.inbounds)
	g.POST("/inbounds", a.getInbounds)
	g.POST("/inbound/add", a.addInbound)
	g.POST("/inbound/del/:id", a.delInbound)
	g.POST("/inbound/update/:id", a.updateInbound)
	g.GET("/setting", a.setting)
	g.POST("/setting/all", a.getAllSetting)
	g.POST("/setting/update", a.updateSetting)
}

func (a *XUIController) startTask() {
	webServer := global.GetWebServer()
	c := webServer.GetCron()
	c.AddFunc("@every 10s", func() {
		if a.isNeedXrayRestart.Load() {
			err := a.xrayService.RestartXray()
			if err != nil {
				logger.Error("restart xray failed:", err)
			}
			a.isNeedXrayRestart.Store(false)
		}
	})
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

func (a *XUIController) getInbounds(c *gin.Context) {
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
	if err == nil {
		a.isNeedXrayRestart.Store(true)
	}
}

func (a *XUIController) delInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "删除", err)
		return
	}
	err = a.inboundService.DelInbound(id)
	jsonMsg(c, "删除", err)
	if err == nil {
		a.isNeedXrayRestart.Store(true)
	}
}

func (a *XUIController) updateInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	inbound := &model.Inbound{
		Id: id,
	}
	err = c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	err = a.inboundService.UpdateInbound(inbound)
	jsonMsg(c, "修改", err)
	if err == nil {
		a.isNeedXrayRestart.Store(true)
	}
}

func (a *XUIController) getAllSetting(c *gin.Context) {
	allSetting, err := a.settingService.GetAllSetting()
	if err != nil {
		jsonMsg(c, "获取设置", err)
		return
	}
	jsonObj(c, allSetting, nil)
}

func (a *XUIController) updateSetting(c *gin.Context) {
	allSetting := &entity.AllSetting{}
	err := c.ShouldBind(allSetting)
	if err != nil {
		jsonMsg(c, "修改设置", err)
		return
	}
	err = a.settingService.UpdateAllSetting(allSetting)
	jsonMsg(c, "修改设置", err)
}
