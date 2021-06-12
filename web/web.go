package web

import (
	"context"
	"crypto/tls"
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/robfig/cron/v3"
	"golang.org/x/text/language"
	"html/template"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"x-ui/config"
	"x-ui/logger"
	"x-ui/util/common"
	"x-ui/web/controller"
	"x-ui/web/service"
)

//go:embed assets/*
var assetsFS embed.FS

//go:embed html/*
var htmlFS embed.FS

//go:embed translation/*
var i18nFS embed.FS

type wrapAssetsFS struct {
	embed.FS
}

func (f *wrapAssetsFS) Open(name string) (fs.File, error) {
	return f.FS.Open("assets/" + name)
}

type Server struct {
	listener net.Listener

	index  *controller.IndexController
	server *controller.ServerController
	xui    *controller.XUIController

	xrayService    service.XrayService
	settingService service.SettingService
	inboundService service.InboundService

	cron *cron.Cron

	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) getHtmlFiles() ([]string, error) {
	files := make([]string, 0)
	dir, _ := os.Getwd()
	err := fs.WalkDir(os.DirFS(dir), "web/html", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *Server) getHtmlTemplate(funcMap template.FuncMap) (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	err := fs.WalkDir(htmlFS, "html", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			newT, err := t.ParseFS(htmlFS, path+"/*.html")
			if err != nil {
				// ignore
				return nil
			}
			t = newT
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Server) initRouter() (*gin.Engine, error) {
	if config.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	secret, err := s.settingService.GetSecret()
	if err != nil {
		return nil, err
	}

	basePath, err := s.settingService.GetBasePath()
	if err != nil {
		return nil, err
	}

	store := cookie.NewStore(secret)
	engine.Use(sessions.Sessions("session", store))
	engine.Use(func(c *gin.Context) {
		c.Set("base_path", basePath)
	})
	err = s.initI18n(engine)
	if err != nil {
		return nil, err
	}

	if config.IsDebug() {
		// for develop
		files, err := s.getHtmlFiles()
		if err != nil {
			return nil, err
		}
		engine.LoadHTMLFiles(files...)
		engine.StaticFS(basePath+"assets", http.FS(os.DirFS("web/assets")))
	} else {
		// for prod
		t, err := s.getHtmlTemplate(engine.FuncMap)
		if err != nil {
			return nil, err
		}
		engine.SetHTMLTemplate(t)
		engine.StaticFS(basePath+"assets", http.FS(&wrapAssetsFS{FS: assetsFS}))
	}

	g := engine.Group(basePath)

	s.index = controller.NewIndexController(g)
	s.server = controller.NewServerController(g)
	s.xui = controller.NewXUIController(g)

	return engine, nil
}

func (s *Server) initI18n(engine *gin.Engine) error {
	bundle := i18n.NewBundle(language.SimplifiedChinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	err := fs.WalkDir(i18nFS, "translation", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := i18nFS.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = bundle.ParseMessageFileBytes(data, path)
		return err
	})
	if err != nil {
		return err
	}

	findI18nParamNames := func(key string) []string {
		names := make([]string, 0)
		keyLen := len(key)
		for i := 0; i < keyLen-1; i++ {
			if key[i:i+2] == "{{" { // 判断开头 "{{"
				j := i + 2
				isFind := false
				for ; j < keyLen-1; j++ {
					if key[j:j+2] == "}}" { // 结尾 "}}"
						isFind = true
						break
					}
				}
				if isFind {
					names = append(names, key[i+3:j])
				}
			}
		}
		return names
	}

	var localizer *i18n.Localizer

	engine.FuncMap["i18n"] = func(key string, params ...string) (string, error) {
		names := findI18nParamNames(key)
		if len(names) != len(params) {
			return "", common.NewError("find names:", names, "---------- params:", params, "---------- num not equal")
		}
		templateData := map[string]interface{}{}
		for i := range names {
			templateData[names[i]] = params[i]
		}
		return localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    key,
			TemplateData: templateData,
		})
	}

	engine.Use(func(c *gin.Context) {
		accept := c.GetHeader("Accept-Language")
		localizer = i18n.NewLocalizer(bundle, accept)
		c.Set("localizer", localizer)
		c.Next()
	})

	return nil
}

func (s *Server) startTask() {
	err := s.xrayService.RestartXray()
	if err != nil {
		logger.Warning("start xray failed:", err)
	}
	var checkTime = 0
	// 每 30 秒检查一次 xray 是否在运行
	s.cron.AddFunc("@every 30s", func() {
		if s.xrayService.IsXrayRunning() {
			checkTime = 0
			return
		}
		checkTime++
		if checkTime < 2 {
			return
		}
		s.xrayService.SetIsNeedRestart(true)
	})

	go func() {
		time.Sleep(time.Second * 5)
		// 每 10 秒统计一次流量，首次启动延迟 5 秒，与重启 xray 的时间错开
		s.cron.AddFunc("@every 10s", func() {
			if !s.xrayService.IsXrayRunning() {
				return
			}
			traffics, err := s.xrayService.GetXrayTraffic()
			if err != nil {
				logger.Warning("get xray traffic failed:", err)
				return
			}
			err = s.inboundService.AddTraffic(traffics)
			if err != nil {
				logger.Warning("add traffic failed:", err)
			}
		})
	}()

	// 每分钟检查一次 inbound 流量超出情况
	s.cron.AddFunc("@every 1m", func() {
		needRestart, err := s.inboundService.DisableInvalidInbounds()
		if err != nil {
			logger.Warning("disable invalid inbounds err:", err)
		} else if needRestart {
			s.xrayService.SetIsNeedRestart(true)
		}
	})
}

func (s *Server) Start() (err error) {
	defer func() {
		if err != nil {
			s.Stop()
		}
	}()

	loc, err := s.settingService.GetTimeLocation()
	if err != nil {
		return err
	}
	s.cron = cron.New(cron.WithLocation(loc), cron.WithSeconds())
	s.cron.Start()

	engine, err := s.initRouter()
	if err != nil {
		return err
	}

	certFile, err := s.settingService.GetCertFile()
	if err != nil {
		return err
	}
	keyFile, err := s.settingService.GetKeyFile()
	if err != nil {
		return err
	}
	listen, err := s.settingService.GetListen()
	if err != nil {
		return err
	}
	port, err := s.settingService.GetPort()
	if err != nil {
		return err
	}
	listenAddr := net.JoinHostPort(listen, strconv.Itoa(port))
	var listener net.Listener
	if certFile != "" || keyFile != "" {
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		c := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		listener, err = tls.Listen("tcp", listenAddr, c)
	} else {
		listener, err = net.Listen("tcp", listenAddr)
	}
	if err != nil {
		return err
	}
	if certFile != "" || keyFile != "" {
		logger.Info("web server run https on", listener.Addr())
	} else {
		logger.Info("web server run http on", listener.Addr())
	}
	s.listener = listener

	s.startTask()

	go engine.RunListener(listener)

	return nil
}

func (s *Server) Stop() error {
	s.cancel()
	s.xrayService.StopXray()
	if s.cron != nil {
		s.cron.Stop()
	}
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) GetCtx() context.Context {
	return s.ctx
}

func (s *Server) GetCron() *cron.Cron {
	return s.cron
}
