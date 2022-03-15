package router

import (
	"context"
	"fmt"
	"log"

	"github.com/NormalReedus/lru-cache-microservice/internal/cache"
	"github.com/NormalReedus/lru-cache-microservice/internal/config"
	"github.com/gofiber/fiber/v2"
)

func Start(conf *config.Config, port string, cache *cache.LRUCache) {
	app := fiber.New()

	// Make cache available in all handlers on ctx.UserContext().Value("cache")
	app.Use(createCtxCacheMiddleware(cache))

	// Will loop through cachable endpoints in config and set route handlers + middleware to handle caching on those routes
	setCachingEndpoints(app, conf)

	// Any non-cache / non-cache-busting requests should just proxy directly to the original API
	setDefaultProxies(app, conf.ApiUrl)

	log.Fatal(app.Listen(fmt.Sprintf(":%v", port)))
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
	cacheCtx := context.WithValue(context.Background(), "cache", cache)

	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(cacheCtx)
		ctx.Next()
		return nil
	}
}
