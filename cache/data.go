package cache

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// NewCacheDataFromJSON takes marshaled json data (from cache) and returns it as hydrated CacheData.
// It is used when reading serialized data (headers and body) from the cache.
func NewCacheDataFromJSON(jsonData []byte) CacheData {
	var data CacheData
	json.Unmarshal(jsonData, &data)

	return data
}

// CacheData is used to represent the headers and body of an API response.
type CacheData struct {
	Headers map[string]string // we don't need to stringify headers
	Body    []byte            `json:"body"`
}

// SetHeaders will add all of the CacheData headers to the fiber context of a route handler.
// This is used when sending cached data to the client to make sure the headers are also the same as the original API response.
func (data *CacheData) SetHeaders(ctx *fiber.Ctx) {
	for key, val := range data.Headers {
		ctx.Set(key, val)
	}
}
