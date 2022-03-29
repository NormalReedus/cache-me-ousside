package config

import "testing"

// Make sure all props are set on the config object properlu
func TestLoadRequiredProps(t *testing.T) {
	config := Load("testdata/test.config.json5")

	if config.Capacity == 0 {
		t.Errorf("Expected config.Capacity to not be 0")
	}

	if config.ApiUrl == "" {
		t.Errorf("Expected config.ApiUrl to not be empty")
	}

	// TODO: change to be a map of []string with GET and HEAD (and more?)
	if len(config.Cache) == 0 {
		t.Errorf("Expected config.Cache to not be empty")
	}

	if len(config.BustMap) == 0 {
		t.Errorf("Expected config.BustMap to not be empty")
	}
	if len(config.BustMap["POST"]["/posts"]) == 0 {
		t.Errorf("Expected config.BustMap's POST /posts endpoint to not be empty")
	}
	if len(config.BustMap["PUT"]["/posts/:slug"]) == 0 {
		t.Errorf("Expected config.BustMap's PUT /posts/:slug endpoint to not be empty")
	}
	if len(config.BustMap["DELETE"]["/posts/:id"]) == 0 {
		t.Errorf("Expected config.BustMap's DELETE /posts/:id endpoint to not be empty")
	}
}

func TestTrimTrailingSlash(t *testing.T) {
	config := Load("testdata/test.config.json5")

	if config.ApiUrl[:len(config.ApiUrl)-1] == "/" {
		t.Errorf("Expected config.ApiUrl to not end with a trailing slash")
	}
}
