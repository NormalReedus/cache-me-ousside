package main

import (
	"fmt"

	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/NormalReedus/cache-me-ousside/internal/router"
)

func main() {
	conf := createConfFromCli()

	logFile := logger.Initialize("log.log")
	if logFile != nil {
		defer logFile.Close()
	}

	fmt.Println(conf) //TODO print the prettified / human-readable configuration

	dataCache := cache.New(conf.Capacity)

	router.Start(conf, "3000", dataCache)
}
