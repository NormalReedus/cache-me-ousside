package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/flynn/json5"
)

var CACHEABLE_METHODS = [...]string{"GET", "HEAD"}
var UNCACHEABLE_METHODS = [...]string{"POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}

func LoadJSON(configPath string) *Config {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		logger.Panic(err)
	}
	defer jsonFile.Close()

	jsonByteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.Panic(err)
	}

	var config = &Config{}
	json5.Unmarshal(jsonByteValue, &config)

	return config
}

// Just has to initialize everything non-primitive so we don't assign to nil-maps
func New() *Config {
	bustMap := make(BustMap)
	bustMap["POST"] = make(map[string][]string)
	bustMap["PUT"] = make(map[string][]string)
	bustMap["DELETE"] = make(map[string][]string)
	bustMap["PATCH"] = make(map[string][]string)
	bustMap["TRACE"] = make(map[string][]string)
	bustMap["CONNECT"] = make(map[string][]string)
	bustMap["OPTIONS"] = make(map[string][]string)

	conf := &Config{
		Cache: make([]string, 0),
		Bust:  bustMap,
	}

	return conf
}

type Config struct {
	Capacity     uint64   `json:"capacity"`     // required
	CapacityUnit string   `json:"capacityUnit"` // Used if you want memory based cache limit
	ApiUrl       string   `json:"apiUrl"`       // required
	Cache        []string `json:"cache"`        // required
	Bust         BustMap  `json:"bust"`
}

// TODO: add support for caching HEAD requests as well
// This requires creating a type for conf.Cache that works like bustMap
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

func (conf *Config) TrimTrailingSlash() {
	conf.ApiUrl = strings.TrimSuffix(conf.ApiUrl, "/")
}

func (conf *Config) ValidateRequiredProps() error {
	if conf.Capacity == 0 {
		return fmt.Errorf("missing required property: capacity")
	}

	if conf.ApiUrl == "" {
		return fmt.Errorf("missing required property: apiUrl")
	}

	// TODO: edit this when cache changes to a map of methods
	if conf.Cache == nil || len(conf.Cache) == 0 {
		return fmt.Errorf("missing required property: cache")
	}

	return nil
}

func (conf Config) String() string {
	// TODO: make this print in a non-json format to display configuration when server runs
	confJSON, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		logger.Warn("there was an issue printing the configuration")
	}
	return string(confJSON)
}

// Is a map of methods with maps of endpoints with slices of patterns to match to cache entries to bust.
type BustMap map[string]map[string][]string
