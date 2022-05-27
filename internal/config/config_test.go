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

	config, _ := LoadJSON(configPath)

	assert.NotZero(config.Capacity, "Expected required prop config.Capacity to be loaded correctly as a non-zero value")
	assert.Positive(config.Capacity, "Expected required prop config.Capacity to be loaded correctly as a positive number")

	assert.NotZero(config.Hostname, "Expected required prop config.Hostname to be loaded correctly as a non-zero value")
	assert.NotZero(config.Port, "Expected required prop config.Port to be loaded correctly as a non-zero value")
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

	conf, err := LoadJSON(configPath)

	assert.Nil(t, conf, "Expected config.LoadJSON to return a nil Config pointer when the config file does not exist")
	assert.Error(t, err, "Expected config.LoadJSON to return an error when the config file does not exist")
}

func TestRequiredProps(t *testing.T) {
	type args struct {
		testFileIdentifier string
		property           string
	}
	tests := [...]args{
		{"api-url-missing", "ApiUrl"},
		{"cache-missing", "Cache"},
		{"cache-empty", "Cache"},
	}

	for _, tt := range tests {
		configPath := "testdata/" + tt.testFileIdentifier + ".json5"

		assert.FileExists(t, configPath, "Expected test configuration file to exist for test to work")

		conf, err := LoadJSON(configPath)

		assert.Nil(t, conf, "Expected config.LoadJSON(\"%s\") to return a nil Config pointer when it does not have required prop: %q", configPath, tt.property)
		assert.Error(t, err, "Expected config.LoadJSON(\"%s\") to return an error when it does not have required prop: %q", configPath, tt.property)

	}
}

func TestTrimTrailingSlash(t *testing.T) {
	configPath := "testdata/test.config.json5"

	assert.FileExists(t, configPath, "Expected test configuration file to exist for test to work")

	conf1 := New()
	conf1.ApiUrl = "https://jsonplaceholder.typicode.com/"
	conf1.TrimTrailingSlash()
	assert.Equal(t, "https://jsonplaceholder.typicode.com", conf1.ApiUrl, "Expected config.TrimTrailingSlash to remove trailing slashes from the api url prop, got: %s", conf1.ApiUrl)

	conf2, _ := LoadJSON(configPath)
	assert.Equal(t, "https://jsonplaceholder.typicode.com", conf2.ApiUrl, "Expected config.LoadJSON to remove trailing slashes from the api url prop when initialized, got: %s", conf2.ApiUrl)

}
