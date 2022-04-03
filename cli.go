package main

import (
	"fmt"
	"time"

	"github.com/NormalReedus/cache-me-ousside/internal/config"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/urfave/cli/v2"
)

// const DEFAULT_CONFIG_PATH = "./cache.config.json5"

// TODO: https://github.com/urfave/cli/blob/master/docs/v2/manual.md#getting-started
// set up with cli, making everything in config file optional
type CLIArgs struct {
	capacity     uint64
	capacityUnit string
	apiUrl       string
	cacheGET     cli.StringSlice // will contain all the paths to cache on GET requests
	cacheHEAD    cli.StringSlice // will contain all the paths to cache on HEAD requests
	bustPOST     cli.StringSlice // first element is the path, rest are the patterns of entries to bust
	bustPUT      cli.StringSlice // first element is the path, rest are the patterns of entries to bust
	bustDELETE   cli.StringSlice // first element is the path, rest are the patterns of entries to bust
	bustPATCH    cli.StringSlice // first element is the path, rest are the patterns of entries to bust
	bustTRACE    cli.StringSlice // first element is the path, rest are the patterns of entries to bust
	bustCONNECT  cli.StringSlice // first element is the path, rest are the patterns of entries to bust
	bustOPTIONS  cli.StringSlice // first element is the path, rest are the patterns of entries to bust
}

func (a *CLIArgs) addToConfig(c *config.Config) {
	if c == nil {
		c = &config.Config{}
	}

	if a.capacity != 0 {
		c.Capacity = a.capacity
	}
	if a.capacityUnit != "" {
		c.CapacityUnit = a.capacityUnit
	}
	fmt.Println("apiUrl:", a.apiUrl)
	if a.apiUrl != "" {
		c.ApiUrl = a.apiUrl
	}
	if len(a.cacheGET.Value()) > 0 {
		c.Cache = a.cacheGET.Value()
	}
	// if len(a.cacheHEAD.Value()) > 0 {
	// 	c.Cache = a.cacheHead.Value()
	// }

	// These must have 2 or more elements, since first element should be the path
	if len(a.bustPOST.Value()) >= 2 {
		c.Bust["POST"][a.bustPOST.Value()[0]] = a.bustPOST.Value()[1:]
	}
	if len(a.bustPUT.Value()) >= 2 {
		c.Bust["PUT"][a.bustPUT.Value()[0]] = a.bustPUT.Value()[1:]
	}
	if len(a.bustDELETE.Value()) >= 2 {
		c.Bust["DELETE"][a.bustDELETE.Value()[0]] = a.bustDELETE.Value()[1:]
	}
	if len(a.bustPATCH.Value()) >= 2 {
		c.Bust["PATCH"][a.bustPATCH.Value()[0]] = a.bustPATCH.Value()[1:]
	}
	if len(a.bustTRACE.Value()) >= 2 {
		c.Bust["TRACE"][a.bustTRACE.Value()[0]] = a.bustTRACE.Value()[1:]
	}
	if len(a.bustCONNECT.Value()) >= 2 {
		c.Bust["CONNECT"][a.bustCONNECT.Value()[0]] = a.bustCONNECT.Value()[1:]
	}
	if len(a.bustOPTIONS.Value()) >= 2 {
		c.Bust["OPTIONS"][a.bustOPTIONS.Value()[0]] = a.bustOPTIONS.Value()[1:]
	}
}

func parseCli() *cli.App {
	args := CLIArgs{}

	app := &cli.App{
		Name:      "cache-me-ousside",
		Version:   "0.0.1",
		Compiled:  time.Now(),
		Copyright: "(c) 2022 Magnus Bendix Borregaard",
		Authors: []*cli.Author{
			{
				Name:  "Magnus Bendix Borregaard",
				Email: "magnus.borregaard@gmail.com",
			},
		},

		Usage: "Sets up an LRU cache microservice that will proxy all your requests to a specified REST API and cache the responses.",

		Flags: []cli.Flag{
			&cli.Uint64Flag{
				Destination: &args.capacity,
				Name:        "capacity",
				Aliases:     []string{"c"},
				Usage:       "the `NUMBER` of entries to cache. If capacity-unit is specfied, this will instead be used as the amount of memory to use for the cache",
			},
			&cli.StringFlag{
				Destination: &args.capacityUnit,
				Name:        "capacity-unit",
				Aliases:     []string{"cu"},
				Usage:       "set this to use a memory-based instead of entry-based cache capacity. Valid `UNIT`s are 'b', 'kb', 'mb', 'gb', and 'tb'",
			},
			&cli.StringFlag{
				Destination: &args.apiUrl,
				Name:        "api-url",
				Aliases:     []string{"u"},
				Usage:       "the `URL` of the API to cache",
			},
			&cli.StringSliceFlag{
				Destination: &args.cacheGET,
				Name:        "cache:GET",
				Aliases:     []string{"c:GET"},
				Usage:       "the list of `PATHS` to cache on GET requests",
			},
			&cli.StringSliceFlag{
				Destination: &args.cacheHEAD,
				Name:        "cache:HEAD",
				Aliases:     []string{"c:HEAD"},
				Usage:       "the list of `PATHS` to cache on HEAD requests",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPOST,
				Name:        "bust:POST",
				Aliases:     []string{"b:POST"},
				Usage:       "first element passed is the path on which a POST request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPUT,
				Name:        "bust:PUT",
				Aliases:     []string{"b:PUT"},
				Usage:       "first element passed is the path on which a PUT request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustDELETE,
				Name:        "bust:DELETE",
				Aliases:     []string{"b:DELETE"},
				Usage:       "first element passed is the path on which a DELETE request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPATCH,
				Name:        "bust:PATCH",
				Aliases:     []string{"b:PATCH"},
				Usage:       "first element passed is the path on which a PATCH request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustTRACE,
				Name:        "bust:TRACE",
				Aliases:     []string{"b:TRACE"},
				Usage:       "first element passed is the path on which a TRACE request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustCONNECT,
				Name:        "bust:CONNECT",
				Aliases:     []string{"b:CONNECT"},
				Usage:       "first element passed is the path on which a CONNECT request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustOPTIONS,
				Name:        "bust:OPTIONS",
				Aliases:     []string{"b:OPTIONS"},
				Usage:       "first element passed is the path on which an OPTIONS request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
		},

		Action: func(c *cli.Context) error {
			var configPath string
			var conf = &config.Config{}

			if c.NArg() > 0 {
				configPath = c.Args().Get(0)
			}

			if configPath != "" {
				conf = config.LoadJSON(configPath)
			}

			args.addToConfig(conf)

			fmt.Printf("%+v", conf)

			if err := conf.ValidateRequiredProps(); err != nil {
				logger.Panic(err)
			}
			conf.TrimTrailingSlash()

			run(conf, "3000")

			return nil
		},
	}

	return app
}
