package cache

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func NewCacheDataFromJSON(jsonData []byte) *CacheData {
	var data CacheData
	json.Unmarshal(jsonData, &data)

	return &data
}

type CacheData struct {
	Headers map[string]string // we don't need to stringify headers
	Body    []byte            `json:"body"`
}

func (data *CacheData) SetHeaders(ctx *fiber.Ctx) {
	for key, val := range data.Headers {
		ctx.Set(key, val)
	}
}
