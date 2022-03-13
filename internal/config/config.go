package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/flynn/json5"
)

func Load(configPath string) *Config {

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonByteValue, _ := ioutil.ReadAll(jsonFile)

	var config = &Config{}

	json5.Unmarshal(jsonByteValue, &config)

	return config
}

type Config struct {
	MaxSize int      `json:"maxSize"`
	ApiUrl  string   `json:"apiUrl"`
	Cache   []string `json:"cache"`
	Bust    BustMap
}

type BustMap struct {
	GET    map[string][]string `json:"GET"`
	HEAD   map[string][]string `json:"HEAD"`
	POST   map[string][]string `json:"POST"`
	PUT    map[string][]string `json:"PUT"`
	DELETE map[string][]string `json:"DELETE"`
	PATCH  map[string][]string `json:"PATCH"`
}
