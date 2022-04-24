package main

import (
	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/NormalReedus/cache-me-ousside/internal/router"
)

const PORT = "3000"

func main() {
	// Initialize logger in terminal mode to log any startup errors to stdout before a potential log file is provided
	logger.Initialize("") // we want all startup errors etc to be logged to terminal, then we will log to file later if one is provided

	// Get configuration struct from CLI (which might read a config file, if provided)
	conf := createConfFromCli()

	// Create the actual cache to hold entries
	dataCache := cache.New(conf.Capacity)

	// Setup the router
	app := router.New(conf, dataCache)

	// Say hello in terminal
	logger.HiMom(conf.ApiUrl, PORT)

	// Set logger to use log file if any is provided
	if conf.LogFilePath != "" {
		logFile := logger.Initialize(conf.LogFilePath)
		if logFile != nil {
			defer logFile.Close()
		}
	}

	// Start the server
	logger.Panic(app.Listen("localhost:" + PORT))
}
