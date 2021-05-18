package main

import (
	"github.com/op/go-logging"
	"log"
	"x-ui/config"
	"x-ui/database"
	"x-ui/logger"
	"x-ui/web"
)

func main() {
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

	server := web.NewServer()
	err = server.Run()
	if err != nil {
		log.Println(err)
	}
}
