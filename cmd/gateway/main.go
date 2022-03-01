package main

import (
	"github.com/def-stocks-t/gateway/internal/config"
	"github.com/def-stocks-t/gateway/internal/rest"
	"github.com/deface90/go-logger/filename"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.New()

	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	logger.AddHook(filenameHook)

	var conf config.Config
	err := configor.Load(&conf, "config.json")
	if err != nil {
		log.Errorf("Failed to read config.json, using default config values")
	}

	restService := rest.NewRestService(conf, logger)
	restService.Run()
}
