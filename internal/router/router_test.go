package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceRouteParams(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		patterns []string
		paramMap map[string]string
		expected []string
	}

	tests := [...]args{
		{
			patterns: []string{
				"^/users/:id",        // replace single param
				"^/users/:id/",       // replace single param with trailing slash
				"^/users/:[a-zA-Z]+", // only match route param syntax / ignore regex
				"^/users/[:id]+",     // precede/ignore regex / this might fuck up regexes
			},
			paramMap: map[string]string{"id": "123"},
			expected: []string{
				"^/users/123",        // replace single param
				"^/users/123/",       // replace single param with trailing slash
				"^/users/:[a-zA-Z]+", // only match route param syntax / ignore regex
				"^/users/[123]+",     // precede/ignore regex / this might fuck up regexes
			},
		},
		{
			patterns: []string{"^/users/:id"}, // don't replace non-existent params in pattern
			paramMap: map[string]string{},
			expected: []string{"^/users/:id"},
		},
		{
			patterns: []string{"^/users/:id"}, // don't do anything with params that do not exist in pattern
			paramMap: map[string]string{"notid": "123"},
			expected: []string{"^/users/:id"},
		},
		{
			patterns: []string{"^/:author/posts/:slug/"}, // replace multiple params
			paramMap: map[string]string{"author": "magnus", "slug": "my-post"},
			expected: []string{"^/magnus/posts/my-post/"},
		},
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, replaceRouteParams(tt.paramMap, tt.patterns))
	}
}
