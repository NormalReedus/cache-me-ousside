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
	missingProps := make([]string, 0)

	if conf.Capacity == 0 {
		missingProps = append(missingProps, "Capacity")
	}

	if conf.ApiUrl == "" {
		missingProps = append(missingProps, "ApiUrl")
	}

	// If cache is missing, empty, or it doesn't have either
	//TODO: add support for just one Cache array that is used for both HEAD and GET
	getExists := conf.cachePropExists("GET")
	headExists := conf.cachePropExists("HEAD")

	if !getExists && !headExists {
		missingProps = append(missingProps, "Cache")
	} else if !getExists {
		missingProps = append(missingProps, "Cache[\"GET\"]")
	} else if !headExists {
		missingProps = append(missingProps, "Cache[\"HEAD\"]")
	}

	if len(missingProps) > 0 {
		return fmt.Errorf("Config is missing the following required properties: %s", strings.Join(missingProps, ", "))
	}

	return nil
}

func (conf *Config) cachePropExists(prop string) bool {
	if conf.Cache[prop] == nil || len(conf.Cache[prop]) == 0 {
		return false
	}

	return true
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
