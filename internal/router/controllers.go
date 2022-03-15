package router

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

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
	//* cache := ctx.UserContext().Value("cache") // how to reference the cache

	fmt.Printf("HIT the cache and return response from cached key %v here, or...\n", ctx.OriginalURL())
	fmt.Println("MISS the cache and Next() to proxy the request")

	ctx.Next()
	return nil
}

//* remember to use cache hit headers etc
func writeCacheMiddleware(ctx *fiber.Ctx) error {
	fmt.Printf("Save the Response() in cache under the key %v here\n", ctx.OriginalURL())

	return nil
}
