package main

import (
	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/NormalReedus/cache-me-ousside/internal/router"
)

func main() {
	conf := createConfFromCli()

	// If there is no logfile set, just use stdout
	logFile := logger.Initialize(conf.LogFilePath)
	if logFile != nil {
		defer logFile.Close()
	}

	dataCache := cache.New(conf.Capacity)

	router.Start(conf, "3000", dataCache)
}
