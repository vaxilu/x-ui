package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"runtime"
	"time"
	"x-ui/web/service"
)

func stopServerController(a *ServerController) {
	a.stopTask()
}

type ServerController struct {
	*serverController
}

func NewServerController(g *gin.RouterGroup) *ServerController {
	a := &ServerController{
		serverController: newServerController(g),
	}
	runtime.SetFinalizer(a, stopServerController)
	return a
}

type serverController struct {
	BaseController

	serverService service.ServerService

	ctx    context.Context
	cancel context.CancelFunc

	lastStatus        *service.Status
	lastGetStatusTime time.Time

	lastVersions        []string
	lastGetVersionsTime time.Time
}

func newServerController(g *gin.RouterGroup) *serverController {
	ctx, cancel := context.WithCancel(context.Background())
	a := &serverController{
		ctx:               ctx,
		cancel:            cancel,
		lastGetStatusTime: time.Now(),
	}
	a.initRouter(g)
	a.startTask()
	return a
}

func (a *serverController) initRouter(g *gin.RouterGroup) {
	g.POST("/server/status", a.status)
	g.POST("/server/getXrayVersion", a.getXrayVersion)
	g.POST("/server/installXray/:version", a.installXray)
}

func (a *serverController) refreshStatus() {
	status := a.serverService.GetStatus(a.lastStatus)
	a.lastStatus = status
}

func (a *serverController) startTask() {
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				break
			default:
			}
			now := time.Now()
			if now.Sub(a.lastGetStatusTime) > time.Minute*3 {
				time.Sleep(time.Second * 2)
				continue
			}
			a.refreshStatus()
		}
	}()
}

func (a *serverController) stopTask() {
	a.cancel()
}

func (a *serverController) status(c *gin.Context) {
	a.lastGetStatusTime = time.Now()

	jsonObj(c, a.lastStatus, nil)
}

func (a *serverController) getXrayVersion(c *gin.Context) {
	now := time.Now()
	if now.Sub(a.lastGetVersionsTime) <= time.Minute {
		jsonObj(c, a.lastVersions, nil)
		return
	}

	versions, err := a.serverService.GetXrayVersions()
	if err != nil {
		jsonMsg(c, "获取版本", err)
		return
	}

	a.lastVersions = versions
	a.lastGetVersionsTime = time.Now()

	jsonObj(c, versions, nil)
}

func (a *serverController) installXray(c *gin.Context) {
	version := c.Param("version")
	err := a.serverService.UpdateXray(version)
	jsonMsg(c, "安装 xray", err)
}
