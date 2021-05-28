package global

import (
	"context"
	"github.com/robfig/cron/v3"
	_ "unsafe"
)

var webServer WebServer

type WebServer interface {
	GetCron() *cron.Cron
	GetCtx() context.Context
}

//go:linkname setWebServer main.setWebServer
func setWebServer(s WebServer) {
	webServer = s
}

func GetWebServer() WebServer {
	return webServer
}
