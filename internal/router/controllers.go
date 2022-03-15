package router

import (
	"github.com/NormalReedus/lru-cache-microservice/internal/cache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// This is used for everything that is not cached
func createProxyHandler(apiUrl string) func(ctx *fiber.Ctx) error {
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

// Decorates createProxyHandler to work as a middleware by also calling Next() after running.
// createProxyHandler must not call Next by itself, since the default handler should always be last.
// this is used for everything that is cached
func createProxyMiddleware(apiUrl string) func(ctx *fiber.Ctx) error {
	proxyHandler := createProxyHandler(apiUrl)

	return func(ctx *fiber.Ctx) error {
		proxyHandler(ctx)

		ctx.Next()
		return nil
	}
}

//* remember to use cache hit headers etc
func readCacheMiddleware(ctx *fiber.Ctx) error {
	cache := ctx.UserContext().Value("cache").(*cache.LRUCache)
	key := ctx.OriginalURL()

	val, present := cache.Get(key)

	// If there is no cached data, continue middlewares to proxy the request
	if !present {
		ctx.Next()
		return nil
	}

	// If there is a key for 'endpoint' in cache.data, send the json from the cache
	//! Client has no idea how to read the []byte val
	//! set headers or save whole request in cache and copy them here?
	ctx.Send(val)

	return nil // don't continue middlewares in this case
}

//* remember to use cache hit headers etc
func writeCacheMiddleware(ctx *fiber.Ctx) error {
	cache := ctx.UserContext().Value("cache").(*cache.LRUCache)

	// Save the response body with cache.Set() and return the result
	key := ctx.OriginalURL()
	body := ctx.Response().Body()

	cache.Set(key, &body)

	return nil // this is always last step, so no Next()
}
