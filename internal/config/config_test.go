package config

import (
	"testing"
)

// Make sure all props are set on the config object properlu
func TestLoadProps(t *testing.T) {
	config := Load("testdata/test.config.json5")

	if config.Capacity == 0 {
		t.Error("Expected config.Capacity to not be 0 when given a valid config file")
	}

	if config.ApiUrl == "" {
		t.Error("Expected config.ApiUrl to not be empty when given a valid config file")
	}

	// TODO: change to be a map of []string with GET and HEAD (and more?)
	if len(config.Cache) == 0 {
		t.Error("Expected config.Cache to not be empty when given a valid config file")
	}

	if len(config.BustMap) == 0 {
		t.Error("Expected config.BustMap to not be empty when given a valid config file")
	}
	if len(config.BustMap["POST"]["/posts"]) == 0 {
		t.Error("Expected config.BustMap's POST /posts endpoint to not be empty when given a valid config file")
	}
	if len(config.BustMap["PUT"]["/posts/:slug"]) == 0 {
		t.Error("Expected config.BustMap's PUT /posts/:slug endpoint to not be empty when given a valid config file")
	}
	if len(config.BustMap["DELETE"]["/posts/:id"]) == 0 {
		t.Error("Expected config.BustMap's DELETE /posts/:id endpoint to not be empty when given a valid config file")
	}
}

func TestBadPathPanic(t *testing.T) {
	defer func() { recover() }()

	configPath := "testdata/does.not.exist.json5"

	Load(configPath)

	t.Errorf("Expected config.Load to panic if the config file: %s doesn't exist", configPath)
}

// TODO: edit this when more props are added to config
func TestRequiredProps(t *testing.T) {
	missingProps := []string{
		"capacity",
		"apiUrl",
		"cacheMap", // TODO: create a version where the prop exists but there are empty slices etc
	}

	for _, prop := range missingProps {
		defer func() { recover() }()

		configPath := "testdata/missing." + prop + ".json5"

		Load(configPath)

		t.Errorf("Expected config.Load to panic when the file: %s is missing the required prop: %s", configPath, prop)
	}
}

func TestTrimTrailingSlash(t *testing.T) {
	config := Load("testdata/test.config.json5")

	if config.ApiUrl[:len(config.ApiUrl)-1] == "/" {
		t.Error("Expected config.Load to remove trailing slashes from the apiUrl")
	}
}
