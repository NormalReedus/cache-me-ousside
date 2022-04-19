package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Make sure all props are set on the config object properlu
func TestLoadProps(t *testing.T) {
	assert := assert.New(t)

	configPath := "testdata/test.config.json5"

	assert.FileExists(configPath, "Expected test configuration file to exist for test to work")

	config := LoadJSON(configPath)

	assert.NotZero(config.Capacity, "Expected required prop config.Capacity to be loaded correctly as a non-zero value")
	assert.Positive(config.Capacity, "Expected required prop config.Capacity to be loaded correctly as a positive number")

	assert.NotZero(config.ApiUrl, "Expected required prop config.ApiUrl to be loaded correctly as a non-zero value")

	// TODO: change to be a map of []string with GET and HEAD (and more?)
	assert.NotEmpty(config.Cache, "Expected config.Cache to not be empty when given a valid config file")

	assert.NotEmpty(config.Bust, "Expected config.Bust to not be empty when given a valid config file")

	assert.NotEmpty(config.Bust["POST"]["/posts"], "Expected config.Bust's POST /posts endpoint to not be empty when given a valid config file")

	assert.NotEmpty(config.Bust["PUT"]["/posts/:slug"], "Expected config.Bust's PUT /posts/:slug endpoint to not be empty when given a valid config file")

	assert.NotEmpty(config.Bust["DELETE"]["/posts/:id"], "Expected config.Bust's DELETE /posts/:id endpoint to not be empty when given a valid config file")
}

func TestBadPathPanic(t *testing.T) {
	configPath := "testdata/does.not.exist.json5"

	assert.NoFileExists(t, configPath, "Expected test configuration file to not exist for test to work")

	assert.Panics(t, func() { LoadJSON(configPath) }, "Expected config.LoadJSON to panic when the config file does not exist")
}

// TODO: edit this when more required props are added to config
func TestRequiredProps(t *testing.T) {
	missingProps := []string{
		"capacity",
		"apiUrl",
		"cache", // TODO: create a version where the prop exists but there are empty slices etc
	}

	for _, prop := range missingProps {
		configPath := "testdata/missing." + prop + ".json5"

		assert.FileExists(t, configPath, "Expected test configuration file to exist for test to work")

		conf := LoadJSON(configPath)

		assert.Error(t, conf.ValidateRequiredProps(), "Expected config.ValidateRequiredProps return an error when the file: %s is missing the required prop: %s", configPath, prop)
	}
}

func TestTrimTrailingSlash(t *testing.T) {
	configPath := "testdata/test.config.json5"

	assert.FileExists(t, configPath, "Expected test configuration file to exist for test to work")

	conf1 := New()
	conf1.ApiUrl = "https://jsonplaceholder.typicode.com/"
	conf1.TrimTrailingSlash()
	assert.Equal(t, "https://jsonplaceholder.typicode.com", conf1.ApiUrl, "Expected config.TrimTrailingSlash to remove trailing slashes from the api url prop, got: %s", conf1.ApiUrl)

	conf2 := LoadJSON(configPath)
	assert.Equal(t, "https://jsonplaceholder.typicode.com", conf2.ApiUrl, "Expected config.LoadJSON to remove trailing slashes from the api url prop when initialized, got: %s", conf2.ApiUrl)

}
