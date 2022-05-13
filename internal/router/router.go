package router

import (
	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/config"
	"github.com/gofiber/fiber/v2"
)

func New(conf *config.Config, cache *cache.LRUCache) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true, // has own HiMom message
		Immutable:             true, // muy importante - makes sure that OriginalUrl() cannot mutate cached endpoints somehow
	})

	// Make cache available in all handlers with ctx.Locals("cache").(*cache.LRUCache)
	app.Use(injectCtxCache(cache))

	// Will loop through methods, endpoints, and patterns and set a middleware for each that removes cache entries when patterns are matched
	setBustingEndpoints(app, conf)

	// Will loop through cachable endpoints in config and set route handlers + middleware to handle caching on those routes
	setCachingEndpoints(app, conf)

	// Any non-cache / non-cache-busting requests should just proxy directly to the original API
	app.Use("*", createProxyHandler(conf.ApiUrl)) // default behavior

	return app
}

// Will loop through methods, endpoints, and patterns and set a middleware for each that removes cache entries when patterns are matched.
// e.g. POST to /users could remove all cache entries that match the pattern ^/users or ^/users/:id etc.
func setBustingEndpoints(app *fiber.App, conf *config.Config) {
	for method, endpointMap := range conf.Bust {
		for endpoint, patterns := range endpointMap {
			app.Add(method, endpoint, createBustMiddleware(patterns))
		}
	}
}

// Sets route handlers and middleware that:
// 1. uses cache if anything is cached, if not, then
// 2. proxies the request to apiUrl and gets a response, then
// 3. sets the response in the cache
func setCachingEndpoints(app *fiber.App, conf *config.Config) {
	// These are the middlewares needed for caching
	middlewares := []func(*fiber.Ctx) error{
		readCacheMiddleware,
		createProxyMiddleware(conf.ApiUrl),
		writeCacheMiddleware,
	}

	// For both cacheable methods, set middlewares on each defined endpoint to cache
	for _, endpoint := range conf.Cache["GET"] {
		app.Get(endpoint, middlewares...)
	}
	for _, endpoint := range conf.Cache["HEAD"] {
		app.Head(endpoint, middlewares...)
	}
}

func injectCtxCache(cache *cache.LRUCache) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("cache", cache)
		ctx.Next()
		return nil
	}
}
