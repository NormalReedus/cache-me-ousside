package commandline

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NormalReedus/cache-me-ousside/internal/config"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/urfave/cli/v2"
)

const ROUTE_SEP_CHAR = "=>"
const PATTERN_SEP_CHAR = "||"

// set up with cli, making everything in config file optional
type CLIArgs struct {
	configPath   string
	capacity     uint64
	capacityUnit string
	hostname     string
	port         uint
	apiUrl       string
	logFilePath  string
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
	if a.hostname != "" {
		c.Hostname = a.hostname
	}
	if a.port != 0 {
		c.Port = a.port
	}
	if a.apiUrl != "" {
		c.ApiUrl = a.apiUrl
	}
	if a.logFilePath != "" {
		c.LogFilePath = a.logFilePath
	}

	if len(a.cacheGET.Value()) > 0 {
		c.Cache["GET"] = a.cacheGET.Value()
	}
	if len(a.cacheHEAD.Value()) > 0 {
		c.Cache["HEAD"] = a.cacheHEAD.Value()
	}

	// Every busting flag should be a string of a path followed by a colon and then a comma separated list of regex patterns to bust
	// so for every Value() we split on colon and comma to set the endpoint and the patterns
	if len(a.bustPOST.Value()) > 0 {
		for _, args := range a.bustPOST.Value() {
			parseAndSetBustArgs(c, "POST", args)
		}
	}
	if len(a.bustPUT.Value()) > 0 {
		for _, args := range a.bustPUT.Value() {
			parseAndSetBustArgs(c, "PUT", args)
		}
	}
	if len(a.bustDELETE.Value()) > 0 {
		for _, args := range a.bustDELETE.Value() {
			parseAndSetBustArgs(c, "DELETE", args)
		}
	}
	if len(a.bustPATCH.Value()) > 0 {
		for _, args := range a.bustPATCH.Value() {
			parseAndSetBustArgs(c, "PATCH", args)
		}
	}
	if len(a.bustTRACE.Value()) > 0 {
		for _, args := range a.bustTRACE.Value() {
			parseAndSetBustArgs(c, "TRACE", args)
		}
	}
	if len(a.bustCONNECT.Value()) > 0 {
		for _, args := range a.bustCONNECT.Value() {
			parseAndSetBustArgs(c, "CONNECT", args)
		}
	}
	if len(a.bustOPTIONS.Value()) > 0 {
		for _, args := range a.bustOPTIONS.Value() {
			parseAndSetBustArgs(c, "OPTIONS", args)
		}
	}

	// Make sure the config is valid
	if err := c.ValidateRequiredProps(); err != nil {
		logger.Panic(err)
	}

	c.TrimTrailingSlash()
	c.RemoveInvalidMethods()
}

func CreateConfFromCli() *config.Config {
	args := CLIArgs{} // holds the flags that should overwrite potential config file values
	var conf *config.Config

	app := &cli.App{
		Name:      "cache-me-ousside",
		Version:   "0.1.0",
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
				EnvVars:     []string{"CONFIG_PATH"},
			},
			&cli.Uint64Flag{
				Destination: &args.capacity,
				Name:        "capacity",
				Aliases:     []string{"cap"},
				Usage:       "the `NUMBER` of entries to cache. If capacity-unit is specfied, this will instead be used as the amount of memory to use for the cache",
				EnvVars:     []string{"CAPACITY"},
			},
			&cli.StringFlag{
				Destination: &args.capacityUnit,
				Name:        "capacity-unit",
				Aliases:     []string{"cu"},
				Usage:       "set this to use a memory-based instead of entry-based cache capacity. Valid `UNIT`s are 'b', 'kb', 'mb', 'gb', and 'tb'",
				EnvVars:     []string{"CAPACITY_UNIT"},
			},
			&cli.StringFlag{
				Destination: &args.hostname,
				Name:        "hostname",
				Aliases:     []string{"host", "hn"},
				Usage:       "the `HOSTNAME` where the cache is accessible",
				EnvVars:     []string{"HOSTNAME"},
			},
			&cli.UintFlag{
				Destination: &args.port,
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "the `PORT` where the cache is accessible",
				EnvVars:     []string{"PORT"},
			},
			&cli.StringFlag{
				Destination: &args.apiUrl,
				Name:        "api-url",
				Aliases:     []string{"url", "u"},
				Usage:       "the `URL` of the API to cache",
				EnvVars:     []string{"PROXY_URL"},
			},
			&cli.StringFlag{
				Destination: &args.logFilePath,
				Name:        "logfile",
				Aliases:     []string{"log", "l"},
				Usage:       "the `FILEPATH` to the log file to use for persistent logs. Omit this to output logs to stdout",
				EnvVars:     []string{"LOGFILE_PATH"},
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
				Usage:       "is parsed from the format '[route]=>[regex-pattern1],[regex-pattern2]...' where regex-patterns are the cache entries to bust when a POST request is made to the route",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPUT,
				Name:        "bust:PUT",
				Aliases:     []string{"b:PUT", "b:put"},
				Usage:       "is parsed from the format '[route]=>[regex-pattern1],[regex-pattern2]...' where regex-patterns are the cache entries to bust when a PUT request is made to the route",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustDELETE,
				Name:        "bust:DELETE",
				Aliases:     []string{"b:DELETE", "b:delete", "b:d"},
				Usage:       "is parsed from the format '[route]=>[regex-pattern1],[regex-pattern2]...' where regex-patterns are the cache entries to bust when a DELETE request is made to the route",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustPATCH,
				Name:        "bust:PATCH",
				Aliases:     []string{"b:PATCH", "b:patch"},
				Usage:       "is parsed from the format '[route]=>[regex-pattern1],[regex-pattern2]...' where regex-patterns are the cache entries to bust when a PATCH request is made to the route",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustTRACE,
				Name:        "bust:TRACE",
				Aliases:     []string{"b:TRACE", "b:trace", "b:t"},
				Usage:       "is parsed from the format '[route]=>[regex-pattern1],[regex-pattern2]...' where regex-patterns are the cache entries to bust when a TRACE request is made to the route",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustCONNECT,
				Name:        "bust:CONNECT",
				Aliases:     []string{"b:CONNECT", "b:connect", "b:c"},
				Usage:       "is parsed from the format '[route]=>[regex-pattern1],[regex-pattern2]...' where regex-patterns are the cache entries to bust when a CONNECT request is made to the route",
			},
			&cli.StringSliceFlag{
				Destination: &args.bustOPTIONS,
				Name:        "bust:OPTIONS",
				Aliases:     []string{"b:OPTIONS", "b:options", "b:o"},
				Usage:       "is parsed from the format '[route]=>[regex-pattern1],[regex-pattern2]...' where regex-patterns are the cache entries to bust when a OPTIONS request is made to the route",
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
			args.addToConfig(conf) // will also trim and validate config

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

func parseAndSetBustArgs(c *config.Config, method, args string) {
	// All busting args must have an arrow (=>) to separate the route from the busting pattern
	if !strings.Contains(args, ROUTE_SEP_CHAR) {
		parseBustArgError(method, args)
	}
	routeAndPatterns := strings.Split(args, ROUTE_SEP_CHAR)
	if len(routeAndPatterns) != 2 {
		parseBustArgError(method, args)
	}

	// First part of the string (before =>) will be the route to listen on
	route := routeAndPatterns[0]
	// Second part of the string (after =>) will be a comma separated list of patterns to bust
	patternsString := routeAndPatterns[1]
	if patternsString == "" {
		parseBustArgError(method, args)
	}

	patterns := strings.Split(routeAndPatterns[1], PATTERN_SEP_CHAR)

	if route == "" || patterns == nil || len(patterns) == 0 {
		parseBustArgError(method, args)
	}

	c.Bust[method][route] = patterns
}

func parseBustArgError(method, args string) {
	logger.Panic(fmt.Errorf("Invalid %s bust argument: %q.\nArgument must be in the format '[route]=>[regex-pattern]||[regex-pattern]...'", method, args))
}
