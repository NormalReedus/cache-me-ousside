package router

import (
	"fmt"
	"log"
	"regexp"

	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// This is used for everything that is not cached
func createProxyHandler(apiUrl string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		url := apiUrl + ctx.OriginalURL()

		if err := proxy.Do(ctx, url); err != nil {
			log.Println(fmt.Errorf("could not proxy request to: %v", url))
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

func readCacheMiddleware(ctx *fiber.Ctx) error {
	dataCache := ctx.Locals("cache").(*cache.LRUCache) // not called 'cache' to avoid conflict with package name
	cacheKey := ctx.OriginalURL()

	cachedData := dataCache.Get(cacheKey)

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

	// Let SysAdmin know they served from cache
	logger.CacheRead(ctx.OriginalURL())

	ctx.Send(cachedData.Body)

	return nil // don't continue middlewares in this case
}

//* remember to use cache hit headers etc
func writeCacheMiddleware(ctx *fiber.Ctx) error {
	// If the response is not a 2xx, don't cache it
	status := ctx.Response().StatusCode()
	if status < 200 && status >= 300 {
		logger.CacheSkip(ctx.OriginalURL())
		return nil
	}

	dataCache := ctx.Locals("cache").(*cache.LRUCache) // not called 'cache' to avoid conflict with package name
	cacheKey := ctx.OriginalURL()

	// Init the current response
	apiResponse := cache.CacheData{
		Headers: ctx.GetRespHeaders(),
		Body:    ctx.Response().Body(),
	}

	// Save the api response in cache
	dataCache.Set(cacheKey, &apiResponse)

	logger.CacheWrite(ctx.OriginalURL())

	return nil // this is always last step, so no Next()
}

func createBustMiddleware(patterns []string) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		dataCache := ctx.Locals("cache").(*cache.LRUCache) // not called 'cache' to avoid conflict with package name

		// Loops through all patterns and finds cache entry keys matching the regex
		var matchedEntries []string
		for _, pattern := range patterns {
			regexp, err := regexp.Compile(pattern)
			if err != nil {
				log.Println(fmt.Errorf("could bust entries with pattern: %s on: %s %s", pattern, ctx.Method(), ctx.OriginalURL()))
			}

			matches := dataCache.Match(regexp)

			if len(matches) > 0 {
				matchedEntries = append(matchedEntries, matches...)
			}
		}

		for _, entryKey := range matchedEntries {
			dataCache.Bust(entryKey)
		}

		ctx.Next()
		return nil
	}
}
