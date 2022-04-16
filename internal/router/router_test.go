package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceRouteParams(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		pattern  string
		paramMap map[string]string
		expected string
	}

	tests := [...]args{
		{pattern: "^/users/:id", paramMap: map[string]string{}, expected: "^/users/:id"},                                                             // don't replace non-existent params in pattern
		{pattern: "^/users/:id", paramMap: map[string]string{"id": "123"}, expected: "^/users/123"},                                                  // replace single param
		{pattern: "^/users/:id/", paramMap: map[string]string{"id": "123"}, expected: "^/users/123/"},                                                // replace single param with trailing slash
		{pattern: "^/users/:id", paramMap: map[string]string{"notid": "123"}, expected: "^/users/:id"},                                               // don't do anything with params that do not exist in pattern
		{pattern: "^/:author/posts/:slug/", paramMap: map[string]string{"author": "magnus", "slug": "my-post"}, expected: "^/magnus/posts/my-post/"}, // replace multiple params
		{pattern: "^/users/:[a-zA-Z]+", paramMap: map[string]string{"id": "123"}, expected: "^/users/:[a-zA-Z]+"},                                    // only match route param syntax / ignore regex
		{pattern: "^/users/[:id]+", paramMap: map[string]string{"id": "123"}, expected: "^/users/[123]+"},                                            // precede/ignore regex / this might fuck up regexes
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, replaceRouteParams(tt.pattern, tt.paramMap))
	}
}
