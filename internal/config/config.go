package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/NormalReedus/cache-me-ousside/cache"
	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/flynn/json5"
	"github.com/olekukonko/tablewriter"
)

var ALL_METHODS = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}
var CACHEABLE_METHODS = ALL_METHODS[0:2]  //  []string{"GET", "HEAD"}
var UNCACHEABLE_METHODS = ALL_METHODS[2:] // []string{"POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}

// Is a map of http methods with maps of endpoints with slices of patterns to match to cache entries to bust.
type BustMap map[string]map[string][]string

// Is a map of http methods with slices of endpoints to cache requests to.
type CacheMap map[string][]string

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
		Cache: make(CacheMap),
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
	config.RemoveInvalidMethods()

	return config
}

type Config struct {
	Capacity     uint64   `json:"capacity"`     // required
	CapacityUnit string   `json:"capacityUnit"` // Used if you want memory based cache limit
	Hostname     string   `json:"hostname"`     // required
	Port         uint     `json:"port"`         // required
	ApiUrl       string   `json:"apiUrl"`       // required
	LogFilePath  string   `json:"logFilePath"`
	Cache        CacheMap `json:"cache"` // required (either GET or HEAD)
	Bust         BustMap  `json:"bust"`
}

// Returns size in bytes or entries as first value and a byteMode bool indicating if the capacity unit is bytes or entries.
// When byteMode is true, the capacity is in bytes, otherwise it is in entries.
func (conf Config) CapacityParsed() (size uint64, byteMode bool) {
	if contains(cache.VALID_CAP_UNITS, strings.ToUpper(conf.CapacityUnit)) {
		bytes, err := cache.ToBytes(conf.Capacity, conf.CapacityUnit)
		if err != nil {
			logger.Panic(err)
		}

		return bytes, true
	}

	return conf.Capacity, false
}

// If capacity is measured in entries, just return the number of entries. Otherwise return the amount of memory the cache can use with the unit appended.
func (conf Config) CapacityString() string {
	cap, byteMode := conf.CapacityParsed()

	if byteMode {
		return fmt.Sprintf("%d%s", conf.Capacity, strings.ToUpper(conf.CapacityUnit)) // e.g., "100 MB"
	}

	return strconv.FormatUint(cap, 10) + " entries" // e.g. "100 entries"
}

// hostname:port
func (conf Config) Address() string {
	return fmt.Sprintf("%s:%d", conf.Hostname, conf.Port)
}

func (conf *Config) TrimTrailingSlash() {
	conf.ApiUrl = strings.TrimSuffix(conf.ApiUrl, "/")
}

// Removes all map keys that are not valid methods and prints a warning to let the user know, that they might have mistyped
func (conf *Config) RemoveInvalidMethods() {
	// Keep track of if an invalid method was spotted, if so, print a list of the valid methods
	invalidCacheMethod := false
	invalidBustMethod := false

	// Only HEAD and GET are valid cacheable methods
	for method := range conf.Cache {
		if !contains(CACHEABLE_METHODS, method) {

			delete(conf.Cache, method)
			logger.Warn(fmt.Sprintf("%s is not a valid cacheable http method, it will be ignored", method))

			invalidCacheMethod = true
		}
	}

	// Bust methods can be any valid method
	for method := range conf.Bust {
		if !contains(ALL_METHODS, method) {

			delete(conf.Bust, method)
			logger.Warn(fmt.Sprintf("%s is not a valid busting http method, it will be ignored", method))

			invalidCacheMethod = true
		}
	}

	if invalidCacheMethod {
		logger.Info(fmt.Sprintf("The following methods are valid cacheable methods: %s", strings.Join(CACHEABLE_METHODS, ", ")))
	}

	if invalidBustMethod {
		logger.Info(fmt.Sprintf("The following methods are valid busting methods: %s", strings.Join(ALL_METHODS, ", ")))
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

	if conf.Hostname == "" {
		missingProps = append(missingProps, "Host")
	}

	if conf.Port == 0 {
		missingProps = append(missingProps, "Port")
	}

	// If cache is missing, empty, or it doesn't have either
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
	// tablewriter writes data to an io.Writer, so we need something that can be written to and converted to a string
	output := new(strings.Builder)

	//* Create general config table
	output.WriteString("\nGeneral Configuration\n")
	generalTable := tablewriter.NewWriter(output)
	generalTable.SetHeader([]string{"Property", "Value"})
	generalTable.SetAutoMergeCells(true)
	generalTable.SetRowLine(true)
	generalTable.AppendBulk([][]string{
		{"Cache address", conf.Address()},
		{"Proxied API URL", conf.ApiUrl},
		{"Capacity", conf.CapacityString()},
		{"Log file", conf.LogFilePath},
	})
	generalTable.Render()

	//* Create cache config table
	output.WriteString("\nCached Endpoints\n")
	cacheRows := [][]string{}
	for _, method := range CACHEABLE_METHODS {
		for _, endpoint := range conf.Cache[method] {
			cacheRows = append(cacheRows, []string{method, endpoint})
		}
	}
	cacheTable := tablewriter.NewWriter(output)
	cacheTable.SetHeader([]string{"Method", "Endpoints"})
	cacheTable.SetAutoMergeCells(true)
	cacheTable.SetRowLine(true)
	cacheTable.AppendBulk(cacheRows)
	cacheTable.Render()

	//* Create bust config table
	output.WriteString("\nBusting Patterns\n")
	bustRows := [][]string{}
	for _, method := range UNCACHEABLE_METHODS {
		for endpoint, endpointMap := range conf.Bust[method] {
			for _, pattern := range endpointMap {
				bustRows = append(bustRows, []string{method, endpoint, pattern})
			}
		}
	}
	bustTable := tablewriter.NewWriter(output)
	bustTable.SetHeader([]string{"Method", "Endpoints", "Patterns"})
	bustTable.SetAutoMergeCells(true)
	bustTable.SetRowLine(true)
	bustTable.AppendBulk(bustRows)
	bustTable.Render()

	return output.String()
}
