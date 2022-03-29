package config

import "testing"

func TestLoadConfig(t *testing.T) {
	config := Load("testdata/test.config.json5")

	if config.ApiUrl[:len(config.ApiUrl)-1] == "/" {
		t.Errorf("Expected config.ApiUrl to not end with a trailing slash")
	}
}
