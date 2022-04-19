package main

import (
	"fmt"
	"os"
	"time"

	"github.com/NormalReedus/cache-me-ousside/internal/config"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/urfave/cli/v2"
)

// const DEFAULT_CONFIG_PATH = "./cache.config.json5"

// TODO: https://github.com/urfave/cli/blob/master/docs/v2/manual.md#getting-started
// set up with cli, making everything in config file optional
type CLIArgs struct {
	configPath   string
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
		c = config.New()
	}

	if a.capacity != 0 {
		c.Capacity = a.capacity
	}
	if a.capacityUnit != "" {
		c.CapacityUnit = a.capacityUnit
	}
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

func createConfFromCli() *config.Config {
	args := CLIArgs{} // holds the flags that should overwrite potential config file values
	var conf *config.Config

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

		Usage:     "Sets up an LRU cache microservice that will proxy all your requests to a specified REST API and cache the responses.",
		ArgsUsage: "first argument passed is an optional json5 config file path",

		Flags: []cli.Flag{
			&cli.PathFlag{
				Destination: &args.configPath,
				Name:        "config",
				Aliases:     []string{"conf", "path"},
				Usage:       "the `PATH` to a json5 config file specifying the cache settings (will be overwritten by command line flags)",
			},
			&cli.Uint64Flag{
				Destination: &args.capacity,
				Name:        "capacity",
				Aliases:     []string{"cap"},
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
				Aliases:     []string{"c:GET", "c:get", "c:g"},
				Usage:       "the list of `PATHS` to cache on GET requests",
			},
			&cli.StringSliceFlag{
				Destination: &args.cacheHEAD,
				Name:        "cache:HEAD",
				Aliases:     []string{"c:HEAD", "c:head", "c:h"},
				Usage:       "the list of `PATHS` to cache on HEAD requests",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPOST,
				Name:        "bust:POST",
				Aliases:     []string{"b:POST", "b:post"},
				Usage:       "first element passed is the path on which a POST request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPUT,
				Name:        "bust:PUT",
				Aliases:     []string{"b:PUT", "b:put"},
				Usage:       "first element passed is the path on which a PUT request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustDELETE,
				Name:        "bust:DELETE",
				Aliases:     []string{"b:DELETE", "b:delete", "b:d"},
				Usage:       "first element passed is the path on which a DELETE request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPATCH,
				Name:        "bust:PATCH",
				Aliases:     []string{"b:PATCH", "b:patch"},
				Usage:       "first element passed is the path on which a PATCH request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustTRACE,
				Name:        "bust:TRACE",
				Aliases:     []string{"b:TRACE", "b:trace", "b:t"},
				Usage:       "first element passed is the path on which a TRACE request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustCONNECT,
				Name:        "bust:CONNECT",
				Aliases:     []string{"b:CONNECT", "b:connect", "b:c"},
				Usage:       "first element passed is the path on which a CONNECT request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustOPTIONS,
				Name:        "bust:OPTIONS",
				Aliases:     []string{"b:OPTIONS", "b:options", "b:o"},
				Usage:       "first element passed is the path on which an OPTIONS request will bust cache entries, subsequent elements are the regex patterns to match to entries to bust",
			},
		},

		Action: func(c *cli.Context) error {
			if c.NArg() > 0 {
				logger.Panic(fmt.Errorf("no arguments should be passed. Did you mean pass a configuration file path with --config?"))
			}

			// If a config path option was passed, initialize config from that file
			if args.configPath != "" {
				conf = config.LoadJSON(args.configPath)
			} else {
				conf = config.New()
			}

			// Add / overwrite cli arguments to config
			args.addToConfig(conf)

			// Make sure the config is valid
			if err := conf.ValidateRequiredProps(); err != nil {
				logger.Panic(err)
			}

			conf.TrimTrailingSlash() // Make sure all routes starting with / will work correctly when proxied
			return nil
		},
	}

	// Use above cli configuration to actually parse cli arguments and create a usable config
	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}

	return conf
}
