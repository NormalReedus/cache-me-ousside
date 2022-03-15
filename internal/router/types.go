package router

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type key string

type CacheResponse struct {
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

func (res *CacheResponse) SetHeaders(ctx *fiber.Ctx) {
	for key, val := range res.Headers {
		ctx.Set(key, val)
	}
}

func NewCacheResponseFromJSON(jsonData []byte) *CacheResponse {
	var res CacheResponse
	json.Unmarshal(jsonData, &res)

	return &res
}

// func (res *ResponseJson) MarshalJSON() ([]byte, error) {

// }
