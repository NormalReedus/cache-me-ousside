package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/flynn/json5"
)

func Load(configPath string) *Config {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonByteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var config = &Config{}

	json5.Unmarshal(jsonByteValue, &config)

	config.trimTrailingSlash()

	return config
}

type Config struct {
	Capacity uint     `json:"capacity"`
	ApiUrl   string   `json:"apiUrl"`
	Cache    []string `json:"cache"`
	Bust     bustMap  `json:"bust"`
}

func (conf *Config) trimTrailingSlash() {
	conf.ApiUrl = strings.TrimSuffix(conf.ApiUrl, "/")
}

func (conf Config) String() string {
	confJSON, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		fmt.Println("there was an issue printing the configuration")
	}
	return string(confJSON)
}

// Is a map of methods with maps of endpoints with slices of patterns to match to cache entries to bust.
type bustMap map[string]map[string][]string
