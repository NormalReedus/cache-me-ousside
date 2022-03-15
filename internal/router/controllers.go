package router

import (
	"encoding/json"
	"fmt"

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
	cache := ctx.UserContext().Value(key("cache")).(*cache.LRUCache)
	cacheKey := ctx.OriginalURL()

	// Read cached json data (headers + body)
	// responseVal is a json with { "headers": { HEADERS JSON }, "body": stringified []byte }
	cachedResponseJson, ok := cache.Get(cacheKey)

	// If there is no cached data, continue middlewares to proxy the request
	if !ok {
		ctx.Set("X-LRU-Cache", "MISS") //! not coming through

		ctx.Next()
		return nil
	}

	// If there is a key for the endpoint in cache, send the json from the cache
	// Init struct with headers and body of cached response
	cachedResponse := NewCacheResponseFromJSON(cachedResponseJson)

	// Set all of the cached headers on the current response
	cachedResponse.SetHeaders(ctx)

	// Let people know they've been served
	ctx.Set("X-LRU-Cache", "HIT")

	ctx.Send(cachedResponse.Body)

	return nil // don't continue middlewares in this case
}

//* remember to use cache hit headers etc
func writeCacheMiddleware(ctx *fiber.Ctx) error {
	cache := ctx.UserContext().Value(key("cache")).(*cache.LRUCache)
	cacheKey := ctx.OriginalURL()

	// Init the current response as a struct we stringify
	apiResponse := CacheResponse{
		Headers: ctx.GetRespHeaders(),
		Body:    ctx.Response().Body(),
	}

	// Stringify the headers + body of the api response
	jsonResponse, err := json.Marshal(apiResponse)
	if err != nil {
		fmt.Println(fmt.Errorf("there was an issue caching the response from %v", cacheKey))
		return nil
	}

	// Save the stringified api response in cache
	cache.Set(cacheKey, &jsonResponse)

	return nil // this is always last step, so no Next()
}
