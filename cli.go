package main

import (
	"github.com/urfave/cli/v2"
)

const DEFAULT_CONFIG_PATH = "./cache.config.json5"

// TODO: https://github.com/urfave/cli/blob/master/docs/v2/manual.md#getting-started
// set up with cli, making everything in config file optional

func Setup() *cli.App {
	app := &cli.App{
		Name: "cache-me-ousside",

		Usage: "Sets up an LRU cache microservice that will proxy all your requests to a specified REST API and cache the responses.",

		Action: func(c *cli.Context) error {
			var configPath string

			if c.NArg() == 0 {
				configPath = DEFAULT_CONFIG_PATH
			} else {
				configPath = c.Args().Get(0)
			}

			run(configPath)
			return nil
		},
	}

	return app
}
