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
	assert.NotZero(config.LogFilePath, "Expected prop config.LogFilePath to be loaded correctly as a non-zero value")

	assert.NotEmpty(config.Cache["HEAD"], "Expected config.Cache[\"HEAD\"] to not be empty when given a valid config file")
	assert.NotEmpty(config.Cache["GET"], "Expected config.Cache[\"GET\"] to not be empty when given a valid config file")

	assert.NotEmpty(config.Bust, "Expected config.Bust to not be empty when given a valid config file")
	assert.NotEmpty(config.Bust["POST"]["/posts"], "Expected config.Bust's POST /posts endpoint to not be empty when given a valid config file")
	assert.NotEmpty(config.Bust["PUT"]["/posts/:id"], "Expected config.Bust's PUT /posts/:id endpoint to not be empty when given a valid config file")
	assert.NotEmpty(config.Bust["DELETE"]["/posts/:id"], "Expected config.Bust's DELETE /posts/:id endpoint to not be empty when given a valid config file")
}

func TestBadPathPanic(t *testing.T) {
	configPath := "testdata/does.not.exist.json5"

	assert.NoFileExists(t, configPath, "Expected test configuration file to not exist for test to work")

	assert.Panics(t, func() { LoadJSON(configPath) }, "Expected config.LoadJSON to panic when the config file does not exist")
}

func TestRequiredProps(t *testing.T) {
	type args struct {
		testFileIdentifier string
		property           string
	}
	tests := [...]args{
		{"capacity-missing", "Capacity"},
		{"api-url-missing", "ApiUrl"},
		{"cache-get-missing", "Cache[\"GET\"]"},
		{"cache-get-empty", "Cache[\"GET\"]"},
		{"cache-head-missing", "Cache[\"HEAD\"]"},
		{"cache-head-empty", "Cache[\"HEAD\"]"},
	}

	for _, tt := range tests {
		configPath := "testdata/" + tt.testFileIdentifier + ".json5"

		assert.FileExists(t, configPath, "Expected test configuration file to exist for test to work")

		assert.Panics(t, func() { LoadJSON(configPath) }, "Expected config.LoadJSON(\"%s\") to panic when it does not have required prop: %q", configPath, tt.property)
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
