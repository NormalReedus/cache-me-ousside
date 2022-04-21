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

var CACHEABLE_METHODS = []string{"GET", "HEAD"}
var ALL_METHODS = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}

// var UNCACHEABLE_METHODS = []string{"POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}

func LoadJSON(configPath string) *Config {
	// Read the configuration json file
	jsonFile, err := os.Open(configPath)
	if err != nil {
		logger.Panic(err)
	}
	defer jsonFile.Close()

	jsonByteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.Panic(err)
	}

	// Populate a new config with the json file values
	var config = New()
	json5.Unmarshal(jsonByteValue, &config)

	// Check if required props are present
	validationErr := config.ValidateRequiredProps()
	if validationErr != nil {
		logger.Panic(validationErr)
	}

	// Clean the API url
	config.TrimTrailingSlash()
	// Remove invalid methods and let the user know
	config.TrimInvalidMethods()

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
		Cache: make(map[string][]string, 0),
		Bust:  bustMap,
	}

	return conf
}

type Config struct {
	Capacity     uint64              `json:"capacity"`     // required
	CapacityUnit string              `json:"capacityUnit"` // Used if you want memory based cache limit
	ApiUrl       string              `json:"apiUrl"`       // required
	Cache        map[string][]string `json:"cache"`        // required
	Bust         BustMap             `json:"bust"`
}

// TODO: add support for caching HEAD requests as well
// This requires creating a type for conf.Cache that works like a map
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

// Removes all map keys that are not valid methods and prints a warning to let the user know, that they might have mistyped
func (conf *Config) TrimInvalidMethods() {
	// Keep track of if an invalid method was spotted, if so, print a list of the valid methods
	invalidCacheMethod := false
	invalidBustMethod := false

	// Only HEAD and GET are valid cacheable methods
	for method := range conf.Cache {
		if !contains(CACHEABLE_METHODS, method) {

			delete(conf.Cache, method)
			logger.Warn(fmt.Sprintf("%q is not a valid cacheable method, it will be ignored", method))

			invalidCacheMethod = true
		}
	}

	// Bust methods can be any valid method
	for method := range conf.Bust {
		if !contains(ALL_METHODS, method) {

			delete(conf.Bust, method)
			logger.Warn(fmt.Sprintf("%q is not a valid busting method, it will be ignored", method))

			invalidCacheMethod = true
		}
	}

	if invalidCacheMethod {
		logger.Info(fmt.Sprintf("The following methods are valid cacheable methods:\n%s", strings.Join(CACHEABLE_METHODS, ", ")))
	}

	if invalidBustMethod {
		logger.Info(fmt.Sprintf("The following methods are valid busting methods:\n%s", strings.Join(ALL_METHODS, ", ")))
	}
}

func (conf *Config) ValidateRequiredProps() error {
	if conf.Capacity == 0 {
		return fmt.Errorf("Config should have a 'Capacity' to know how many entries the cache should hold")
	}

	if conf.ApiUrl == "" {
		return fmt.Errorf("Config should have an 'ApiUrl' to know where to proxy requests to and cache the data from")
	}

	// If cache is missing, empty, or it doesn't have either
	headMissing := conf.Cache["HEAD"] == nil || len(conf.Cache["HEAD"]) == 0
	getMissing := conf.Cache["GET"] == nil || len(conf.Cache["GET"]) == 0
	if conf.Cache == nil || (headMissing && getMissing) {
		return fmt.Errorf("Config should have a list of caching endpoints with their respective HTTP request methods in either 'Cache[\"HEAD\"] or Cache[\"GET\"]")
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
