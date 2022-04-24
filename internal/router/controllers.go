package router

import (
	"fmt"

	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// This is used for everything that is not cached, and as such it does not call Next()
func createProxyHandler(apiUrl string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		url := apiUrl + ctx.OriginalURL()

		if err := proxy.Do(ctx, url); err != nil {
			logger.Error(fmt.Errorf("could not proxy request to: %v", url))
			return err
		}

		// Remove Server header from response
		ctx.Response().Header.Del(fiber.HeaderServer)

		return nil
	}
}

// Decorates createProxyHandler to work as a middleware by also calling Next() after running.
// createProxyHandler must not call Next by itself, since the default handler should always be last.
// This is used for every route that is cached
func createProxyMiddleware(apiUrl string) func(ctx *fiber.Ctx) error {
	proxyHandler := createProxyHandler(apiUrl)

	return func(ctx *fiber.Ctx) error {
		proxyHandler(ctx)

		ctx.Next()
		return nil
	}
}

// This checks for existing entries on a given route and method and returns the cached item if it exists.
// Otherwise we call Next() to proxy the request and get data from the api.
func readCacheMiddleware(ctx *fiber.Ctx) error {
	dataCache := ctx.Locals("cache").(*cache.LRUCache) // not called 'cache' to avoid conflict with package name
	entryKey := generateEntryKey(ctx)                  // which entry to look for in the cache

	cachedData := dataCache.Get(entryKey)

	// If there is no cached data, continue middlewares to proxy the request
	if cachedData == nil {
		ctx.Set("X-LRU-Cache", "MISS")

		ctx.Next()
		return nil
	}

	// Set all of the cached headers on the current response
	cachedData.SetHeaders(ctx)

	// Let people know they've been served
	ctx.Set("X-LRU-Cache", "HIT")

	// Let SysAdmin know they served something from cache
	logger.CacheRead(entryKey)

	ctx.Send(cachedData.Body)

	return nil // don't continue middlewares in this case
}

func writeCacheMiddleware(ctx *fiber.Ctx) error {
	entryKey := generateEntryKey(ctx) // the name to use when saving the entry in cache
	// If the response is not a 2xx, don't cache it
	status := ctx.Response().StatusCode()
	if status < 200 && status >= 300 {
		logger.CacheSkip(entryKey)
		return nil
	}

	dataCache := ctx.Locals("cache").(*cache.LRUCache) // not called 'cache' to avoid conflict with package name

	// Init the current response
	apiResponse := cache.CacheData{
		Headers: ctx.GetRespHeaders(),
		Body:    ctx.Response().Body(),
	}

	// Save the api response in cache
	dataCache.Set(entryKey, &apiResponse)

	logger.CacheWrite(entryKey)

	return nil // this is always last step, so no Next()
}

func createBustMiddleware(patterns []string) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		dataCache := ctx.Locals("cache").(*cache.LRUCache) // not called 'cache' to avoid conflict with package name

		// Now find all cache entries that match the regex pattern or specific route with param
		matchedEntries := dataCache.Match(patterns, ctx.AllParams())

		// Remove the matched entries from the cache
		dataCache.Bust(matchedEntries...)

		ctx.Next()
		return nil
	}
}

func generateEntryKey(ctx *fiber.Ctx) string {
	return ctx.Method() + ":" + ctx.OriginalURL()
}
