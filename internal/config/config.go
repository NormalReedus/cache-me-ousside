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
	Capacity     uint64   `json:"capacity"`
	CapacityUnit string   `json:"capacityUnit"` // Used if you want memory based cache limit
	ApiUrl       string   `json:"apiUrl"`
	Cache        []string `json:"cache"`
	BustMap      bustMap  `json:"bust"`
}

// TODO: add support for caching HEAD requests as well
// Also add support for any other methods that function like GET
// This requires renaming conf.Cache to conf.CacheMap which works kinda like BustMap
// and updating the middleware factory for caching

// TODO: add memory based cache limit
// Create a method to return the cache capacity
// if CapacityUnit is set, use utils.ToBytes to convert the capacity to bytes
// otherwise return the capacity as a number of entries
// (maybe there should be something that tells whether we use entries og memory)
// when busting a cache entry, we should then use utils.MemUsage to compare with the capacity
// when deciding whether to evict, instead of using entries. Using one over the other should
// be checked with a bool on the config that is initialized in the factory function, so busting
// knows whether to use memory or entries

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
