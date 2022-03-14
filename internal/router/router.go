package router

import (
	"fmt"
	"log"

	"github.com/NormalReedus/lru-cache-microservice/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func Start(conf *config.Config, port string) {

	app := fiber.New()

	// Will loop through cachable endpoints in config and set route handlers + middleware to handle caching on those routes
	setCacheMiddleware(app, conf.ApiUrl)

	// Any non-cache / non-cache-busting requests should just proxy directly to the original API
	setDefaultProxies(app, conf.ApiUrl)

	log.Fatal(app.Listen(fmt.Sprintf(":%v", port)))
}

// Sets route handlers and middleware that:
// 1. uses cache if anything is cached, if not, then...
// 2. proxies the request to apiUrl and gets a response, then...
// 3. sets the response in the cache
func setCacheMiddleware(app *fiber.App, apiUrl string) {
	// remember to use cache hit headers etc

	// Set these for all cachable get request routes in Config
	app.Get(
		"/todos/*", // get from config in a loop

		// create service for this
		func(ctx *fiber.Ctx) error {
			fmt.Println("HIT the cache and return response here, or...")
			fmt.Println("MISS the cache and Next() to proxy the request")

			ctx.Next()
			return nil
		},

		createProxyMiddleware(apiUrl),

		// create service for this
		func(ctx *fiber.Ctx) error {
			fmt.Println("Save the Response() in cache here")

			return nil
		},
	)
}

func setDefaultProxies(app *fiber.App, apiUrl string) {
	app.Get("/*", createDefaultProxyHandler(apiUrl))
	app.Head("/*", createDefaultProxyHandler(apiUrl))
	app.Post("/*", createDefaultProxyHandler(apiUrl))
	app.Put("/*", createDefaultProxyHandler(apiUrl))
	app.Patch("/*", createDefaultProxyHandler(apiUrl))
	app.Delete("/*", createDefaultProxyHandler(apiUrl))
}

func createDefaultProxyHandler(apiUrl string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		url := apiUrl + ctx.OriginalURL()

		if err := proxy.Do(ctx, url); err != nil {
			return err
		}

		// Remove Server header from response
		ctx.Response().Header.Del(fiber.HeaderServer)

		return nil
	}
}

// Decorates createDefaultProxyHandler to work as a middleware by also calling Next() after running.
// createDefaultProxyHandler must not call Next by itself, since the default handler should always be last.
func createProxyMiddleware(apiUrl string) func(ctx *fiber.Ctx) error {
	proxyHandler := createDefaultProxyHandler(apiUrl)

	return func(ctx *fiber.Ctx) error {
		proxyHandler(ctx)

		ctx.Next()
		return nil
	}
}
