package main

import (
	"os"

	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/config"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/NormalReedus/cache-me-ousside/internal/router"
)

func main() {
	app := parseCli()

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}

	// var configPath, port string = parseArgs()

}

func run(conf *config.Config, port string) {
	dataCache := cache.New(conf.Capacity)

	router.Start(conf, port, dataCache)
}

// func parseArgs() (string, string) {
// 	helpPtr := flag.BoolP("help", "h", false, "Print help.")
// 	portPtr := flag.StringP("port", "p", "3000", "The port to serve the service on.")
// 	flag.Parse()

// 	if *helpPtr {
// 		printHelp()
// 		os.Exit(0)
// 	}

// 	// Grab config file path from first arg given through CLI
// 	configPath := flag.Arg(0)

// 	// If no config path is given from CLI, use default
// 	if configPath == "" {
// 		configPath = DEFAULT_CONFIG_PATH
// 	}

// 	return configPath, *portPtr
// }

// func printHelp() {
// 	fmt.Print("lru-cache-microservice is a reverse proxy for caching simple requests to a REST API. You only have to configure your API to trust this proxy for optimal conditions.\n\n")
// 	fmt.Println("When running lru-cache-microservice, you will need to supply the program with a configuration file (JSON/JSON5) that specifies which requests to cache, and when to bust the cache.")
// 	fmt.Println("You can find detailed documentation on how to use lru-cache-microservice at https://github.com/NormalReedus/cache-me-ousside/blob/main/README.md.")
// }
