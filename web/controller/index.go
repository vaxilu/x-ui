package controller

import (
	"net/http"
	"time"
	"x-ui/logger"
	"x-ui/web/job"
	"x-ui/web/service"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type IndexController struct {
	BaseController

	userService service.UserService
}

func NewIndexController(g *gin.RouterGroup) *IndexController {
	a := &IndexController{}
	a.initRouter(g)
	return a
}

func (a *IndexController) initRouter(g *gin.RouterGroup) {
	g.GET("/", a.index)
	g.POST("/login", a.login)
	g.GET("/logout", a.logout)
}

func (a *IndexController) index(c *gin.Context) {
	if session.IsLogin(c) {
		c.Redirect(http.StatusTemporaryRedirect, "xui/")
		return
	}
	html(c, "login.html", "Conecte-se", nil)
}

func (a *IndexController) login(c *gin.Context) {
	var form LoginForm
	err := c.ShouldBind(&form)
	if err != nil {
		pureJsonMsg(c, false, "erro de formato de dados")
		return
	}
	if form.Username == "" {
		pureJsonMsg(c, false, "por favor, digite o nome de usuário")
		return
	}
	if form.Password == "" {
		pureJsonMsg(c, false, "Por favor, digite a senha")
		return
	}
	user := a.userService.CheckUser(form.Username, form.Password)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	if user == nil {
		job.NewStatsNotifyJob().UserLoginNotify(form.Username, getRemoteIp(c), timeStr, 0)
		logger.Infof("nome de usuário ou senha errado: \"%s\" \"%s\"", form.Username, form.Password)
		pureJsonMsg(c, false, "nome de usuário ou senha errados")
		return
	} else {
		logger.Infof("%s coenctado,Ip Address:%s\n", form.Username, getRemoteIp(c))
		job.NewStatsNotifyJob().UserLoginNotify(form.Username, getRemoteIp(c), timeStr, 1)
	}

	err = session.SetLoginUser(c, user)
	logger.Info("usuario", user.Id, "coenctado")
	jsonMsg(c, "Conecte-se", err)
}

func (a *IndexController) logout(c *gin.Context) {
	user := session.GetLoginUser(c)
	if user != nil {
		logger.Info("usuario", user.Id, "sair")
	}
	session.ClearSession(c)
	c.Redirect(http.StatusTemporaryRedirect, c.GetString("base_path"))
}
