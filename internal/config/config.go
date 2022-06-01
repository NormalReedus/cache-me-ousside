// Package config loads, validates, and manages the configuration for the LRUCache.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/magnus-bb/cache-me-ousside/cache"
	"github.com/magnus-bb/cache-me-ousside/internal/logger"
	"github.com/olekukonko/tablewriter"
)

const (
	DefaultCapacity uint64 = 500
	DefaultHostname string = "localhost"
	DefaultPort     uint   = 8080
)

var (
	// AllHTTPMethods is a slice of all valid http methods to use in the cache configuration for busting.
	// 	{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}
	AllHTTPMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}

	// CacheableHTTPMethods is a slice of all cacheable http methods that can be used in the cache configuration for caching.
	// 	{"GET", "HEAD"}
	CacheableHTTPMethods = AllHTTPMethods[0:2]

	// UncacheableHTTPMethods = AllMethods[2:] // []string{"POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT", "OPTIONS"}
)

type (
	// BustMap represents a map of http methods with maps of endpoints with slices of patterns to match to cache entries to bust.
	BustMap map[string]map[string][]string
	// CacheMap represents a map of http methods with slices of endpoints to which requests should be cached.
	CacheMap map[string][]string
)

// New returns a Config where Bust and Cache are initialized to empty BustMap and CacheMap respectively.
// This is done to avoid nil pointers when accessing the nested map properties.
func New() *Config {
	bustMap := make(BustMap)
	for _, method := range AllHTTPMethods {
		bustMap[method] = make(map[string][]string)
	}

	conf := &Config{
		Capacity: DefaultCapacity,
		Hostname: DefaultHostname,
		Port:     DefaultPort,
		Cache:    make(CacheMap),
		Bust:     bustMap,
	}

	return conf
}

/*
LoadJSON returns a Config created from unmarshaling the json file at configPath.
It will also trim invalid http methods in the configuration as well as trailing slashes in the ApiUrl.
Be aware that you will manually have to validate the configuration like so:
	conf, _ := LoadJSON(configPath)
	if err := conf.ValidateProps(); err != nil {
		// Handle validation error
	}
*/
func LoadJSON(configPath string) (*Config, error) {
	// Read the configuration json file
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonByteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	// Populate a new config with the json file values
	var config = New()
	json.Unmarshal(jsonByteValue, &config)

	// Clean the API url
	config.TrimTrailingSlash()
	// Remove invalid methods and let the user know
	config.RemoveInvalidHTTPMethods()

	return config, nil
}

// Config represents the configuration for the cache-me-ousside application.
type Config struct {
	// Default is 500, it represents the limit to how much data can be stored in the cache.
	Capacity uint64 `json:"capacity" validate:"min=1"`

	/*
		CapacityUnit represents the unit of measurement for the capacity.
		If omitted, the cache Capacity will be measured in entries.
		Set CapacityUnit to a string to use memory as the unit of measurement, e.g. "mb".
	*/
	CapacityUnit string `json:"capacityUnit" validate:"required,omitempty,ascii,oneof=b B kb KB mb MB gb GB tb TB"`

	// Default is "localhost", it represents the hostname where the server application can be accessed.
	Hostname string `json:"hostname" validate:"required,omitempty,hostname_rfc1123"`

	//Default is 8080, it represents the port where the server application can be accessed. E.g.:
	Port uint `json:"port" validate:"required,min=1,max=65535"`

	// ApiUrl is required, it represents the url of the API to which all requests are proxied and cached from.
	ApiUrl string `json:"apiUrl" validate:"required,url"`

	// LogFilePath is the path to an optional log file to use instead of stdout (terminal mode).
	LogFilePath string `json:"logFilePath" validate:"omitempty,filepath"`

	/*
		Cache is a map of HTTP methods with slices of endpoints to which requests should be cached. E.g.:
			{
				"GET": ["/api/v1/users/:id", "/api/v1/users/:id/posts"],
				"HEAD": ["/api/v1/users/:id", "/api/v1/users/:id/posts"],
			}
	*/
	Cache CacheMap `json:"cache" validate:"required,gt=0,dive,keys,oneof=GET HEAD,endkeys,dive,route"`

	/*
		Bust is a map of HTTP methods with maps of endpoints with slices of patterns to match to cache entries to bust. E.g.:
			{
				"POST": {
					"/posts": [ "/posts" ]
				},
				"PUT": {
					"/posts": [ "^GET:/posts", "^HEAD:/posts" ],
					"/posts/:id": [ "/posts/:id" ]
				}
			}
	*/
	Bust BustMap `json:"bust" validate:"omitempty,dive,keys,oneof=GET HEAD POST PUT DELETE PATCH TRACE CONNECT OPTIONS,endkeys,dive,keys,route"`
}

/*
	CapacityParsed returns size in bytes or entries as first value and a byteMode bool
	indicating if the capacity unit is bytes or entries.
	If CapacityUnit is a valid memory size string, the size is converted from the memory unit to bytes (e.g., mb -> bytes).
*/
func (conf Config) CapacityParsed() (size uint64, byteMode bool) {
	if contains(cache.ValidCapacityUnits, strings.ToUpper(conf.CapacityUnit)) {
		bytes, _ := cache.ToBytes(conf.Capacity, conf.CapacityUnit) // expect memory unit to be valid if the config is validated on creation

		return bytes, true
	}

	return conf.Capacity, false
}

// CapacityString returns a human-readable string representation of the cache capacity.
// If capacity is measured in entries, just return the number of entries. Otherwise return the amount of memory the cache can use with the unit appended.
func (conf Config) CapacityString() string {
	cap, byteMode := conf.CapacityParsed()

	if byteMode {
		return fmt.Sprintf("%d%s", conf.Capacity, strings.ToUpper(conf.CapacityUnit)) // e.g., "100 MB"
	}

	return strconv.FormatUint(cap, 10) + " entries" // e.g. "100 entries"
}

// Address returns the server address in the format hostname:port.
// This is where the server application can be accessed.
func (conf Config) Address() string {
	return fmt.Sprintf("%s:%d", conf.Hostname, conf.Port)
}

// LogModeString returns a human-readable string representation of how logging is configured.
// It will be either a log file path or "terminal mode"
func (conf Config) LogModeString() string {
	if conf.LogFilePath != "" {
		return conf.LogFilePath
	}

	return "terminal mode"
}

// TrimTrailingSlash mutates the ApiUrl to remove any trailing slashes.
// This is useful so all specified endpoints and patterns can begin with a slash.
func (conf *Config) TrimTrailingSlash() {
	conf.ApiUrl = strings.TrimSuffix(conf.ApiUrl, "/")
}

// RemoveInvalidHTTPMethods removes all map keys that are not valid http methods
// and prints a warning to let the user know, that they might have mistyped the configuration.
func (conf *Config) RemoveInvalidHTTPMethods() {
	// Keep track of if an invalid method was spotted, if so, print a list of the valid methods
	invalidCacheMethod := false
	invalidBustMethod := false

	// Only HEAD and GET are valid cacheable methods
	for method := range conf.Cache {
		if !contains(CacheableHTTPMethods, method) {

			delete(conf.Cache, method)
			logger.Warn(fmt.Sprintf("%s is not a valid cacheable HTTP method, it will be ignored", method))

			invalidCacheMethod = true
		}
	}

	// Bust methods can be any valid method
	for method := range conf.Bust {
		if !contains(AllHTTPMethods, method) {

			delete(conf.Bust, method)
			logger.Warn(fmt.Sprintf("%s is not a valid busting HTTP method, it will be ignored", method))

			invalidCacheMethod = true
		}
	}

	if invalidCacheMethod {
		logger.Info("the following methods are valid cacheable HTTP methods: " + strings.Join(CacheableHTTPMethods, ", "))
	}

	if invalidBustMethod {
		logger.Info("the following methods are valid busting HTTP methods: " + strings.Join(AllHTTPMethods, ", "))
	}
}

// ValidateProps makes sure required configuration props are set and follow the correct format.
// TODO: use https://github.com/go-playground/validator.
func (conf *Config) ValidateProps() error {
	validate := newConfigValidator()
	errs := validate.Struct(conf)

	var formattedError error

	if errs != nil {
		formattedError = formatValidationError(errs.(validator.ValidationErrors))
		fmt.Println(formattedError)
	}

	var errorStr string

	if _, err := cache.ToBytes(conf.Capacity, conf.CapacityUnit); err != nil && conf.CapacityUnit != "" {
		errorStr += fmt.Sprintf("%s\n", err)
	}

	missingProps := make([]string, 0)
	if conf.Capacity == 0 {
		missingProps = append(missingProps, "Capacity")
	}

	if conf.ApiUrl == "" {
		missingProps = append(missingProps, "ApiUrl")
	}

	if conf.Hostname == "" {
		missingProps = append(missingProps, "Hostname")
	}

	if conf.Port == 0 {
		missingProps = append(missingProps, "Port")
	}

	// If cache is missing, empty, or it doesn't have either
	getExists := conf.cachePropExists("GET")
	headExists := conf.cachePropExists("HEAD")

	if !getExists && !headExists {
		missingProps = append(missingProps, "Cache")
	}

	if len(missingProps) > 0 {
		errorStr += fmt.Sprintf("Config is missing the following required properties: %s\n", strings.Join(missingProps, ", "))
	}

	if errorStr != "" {
		return errors.New(errorStr)
	}

	return nil
}

// cachePropExists returns true if the cache map has the given prop.
func (conf *Config) cachePropExists(prop string) bool {
	if conf.Cache[prop] == nil || len(conf.Cache[prop]) == 0 {
		return false
	}

	return true
}

// String returns a human-readable table-formatted representation of the configuration.
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
		{"Log", conf.LogModeString()},
	})
	generalTable.Render()

	//* Create cache config table
	output.WriteString("\nCached Endpoints\n")
	cacheRows := [][]string{}
	for _, method := range CacheableHTTPMethods {
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
	bustRows := [][]string{}

	for _, method := range AllHTTPMethods {
		for endpoint, endpointMap := range conf.Bust[method] {
			if len(endpointMap) == 0 {
				bustRows = append(bustRows, []string{method, endpoint, "ALL"}) // empty bust pattern slice means to bust everything
				continue
			}

			for _, pattern := range endpointMap {
				bustRows = append(bustRows, []string{method, endpoint, pattern})
			}
		}
	}

	// No need to print anything, if there are no bust methods declared
	if len(bustRows) != 0 {
		output.WriteString("\nCache Busting Patterns\n")
		bustTable := tablewriter.NewWriter(output)
		bustTable.SetHeader([]string{"Method", "Endpoints", "Patterns"})
		bustTable.SetAutoMergeCells(true)
		bustTable.SetRowLine(true)
		bustTable.AppendBulk(bustRows)
		bustTable.Render()
	}

	return output.String()
}
