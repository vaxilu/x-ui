package main

import (
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "unsafe"
	"x-ui/config"
	"x-ui/database"
	"x-ui/logger"
	"x-ui/web"
	"x-ui/web/global"
)

// this function call global.setWebServer
func setWebServer(server global.WebServer)

func runWebServer() {
	log.Printf("%v %v", config.GetName(), config.GetVersion())

	switch config.GetLogLevel() {
	case config.Debug:
		logger.InitLogger(logging.DEBUG)
	case config.Info:
		logger.InitLogger(logging.INFO)
	case config.Warn:
		logger.InitLogger(logging.WARNING)
	case config.Error:
		logger.InitLogger(logging.ERROR)
	default:
		log.Fatal("unknown log level:", config.GetLogLevel())
	}

	err := database.InitDB(config.GetDBPath())
	if err != nil {
		log.Fatal(err)
	}

	var server *web.Server

	server = web.NewServer()
	setWebServer(server)
	err = server.Start()
	if err != nil {
		panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP)
	for {
		sig := <-sigCh

		if sig == syscall.SIGHUP {
			server.Stop()
			server = web.NewServer()
			setWebServer(server)
			err = server.Start()
			if err != nil {
				panic(err)
			}
		} else {
			continue
		}
	}
}

func v2ui(dbPath string) {
	// migrate from v2-ui
}

func main() {
	if len(os.Args) < 2 {
		runWebServer()
		return
	}

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	v2uiCmd := flag.NewFlagSet("v2-ui", flag.ExitOnError)
	var dbPath string
	v2uiCmd.StringVar(&dbPath, "db", "/etc/v2-ui/v2-ui.db", "set v2-ui db file path")

	switch flag.Arg(0) {
	case "run":
		runCmd.Parse(os.Args[2:])
		runWebServer()
	case "v2-ui":
		v2uiCmd.Parse(os.Args[2:])
		v2ui(dbPath)
	default:
		fmt.Println("excepted 'run' or 'v2-ui' subcommands")
		fmt.Println()
		runCmd.Usage()
		fmt.Println()
		v2uiCmd.Usage()
	}
}
