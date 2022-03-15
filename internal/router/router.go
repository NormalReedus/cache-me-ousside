package router

import (
	"context"
	"fmt"
	"log"

	"github.com/NormalReedus/lru-cache-microservice/internal/cache"
	"github.com/NormalReedus/lru-cache-microservice/internal/config"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

const CACHE_KEY key = "cache"

func Start(conf *config.Config, port string, cache *cache.LRUCache) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Make cache available in all handlers on ctx.UserContext().Value(key("cache"))
	app.Use(createCtxCacheMiddleware(cache))

	// Will loop through cachable endpoints in config and set route handlers + middleware to handle caching on those routes
	setCachingEndpoints(app, conf)

	// Any non-cache / non-cache-busting requests should just proxy directly to the original API
	setDefaultProxies(app, conf.ApiUrl)

	printHiMom(conf, port)

	log.Fatal(app.Listen(fmt.Sprintf("localhost:%v", port)))
}

// Sets route handlers and middleware that:
// 1. uses cache if anything is cached, if not, then...
// 2. proxies the request to apiUrl and gets a response, then...
// 3. sets the response in the cache
func setCachingEndpoints(app *fiber.App, conf *config.Config) {
	for _, endpoint := range conf.Cache {
		app.Get(
			endpoint,

			readCacheMiddleware,

			createProxyMiddleware(conf.ApiUrl), // needs ApiUrl to proxy request

			writeCacheMiddleware,
		)
	}
}

func setDefaultProxies(app *fiber.App, apiUrl string) {
	app.Get("/*", createProxyHandler(apiUrl))
	app.Head("/*", createProxyHandler(apiUrl))
	app.Post("/*", createProxyHandler(apiUrl))
	app.Put("/*", createProxyHandler(apiUrl))
	app.Patch("/*", createProxyHandler(apiUrl))
	app.Delete("/*", createProxyHandler(apiUrl))
}

func createCtxCacheMiddleware(cache *cache.LRUCache) func(ctx *fiber.Ctx) error {
	cacheCtx := context.WithValue(context.Background(), CACHE_KEY, cache)

	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(cacheCtx)
		ctx.Next()
		return nil
	}
}

func printHiMom(conf *config.Config, port string) {
	cacheColor := color.New(color.FgBlue, color.Bold)
	urlColor := color.New(color.FgHiGreen).Add(color.Underline)

	fmt.Print("Your ")
	cacheColor.Print("LRU cache microservice ")
	fmt.Printf("is being served on port %v.\n", port)
	fmt.Print("All requests will be proxied to ")
	urlColor.Print(conf.ApiUrl)
}
