package cache

import (
	"fmt"
	"runtime"
	"strings"
)

//* SET
// Used for finding unique cache entry keys matched by bust-patterns
// (multiple patterns can match same key, we don't need more than one to bust it from cache)

type void struct{}

var nothing void

type Set[T comparable] map[T]void

func (set Set[T]) Add(elem T) {
	set[elem] = nothing
}

// Actually not needed for this specific project
func (set Set[T]) Remove(elem T) {
	delete(set, elem)
}

// Actually not needed for this specific project
func (set Set[T]) Has(elem T) bool {
	_, ok := set[elem]
	return ok
}

func (set Set[T]) Elements() []T {
	elements := make([]T, 0, len(set))
	for k := range set {
		elements = append(elements, k)
	}
	return elements
}

// Takes route patterns (strings to turn into regex) and a map of all route params in a route handler (ctx.AllParams()) and returns the routePattern with route params replaced with their arguments.
// Example: will replace /users/:id with /users/123 when given map[id:123]
func hydrateParams(paramMap map[string]string, routePatternTemplates []string) []string {
	// Copy original slice so we can return a new one
	// newRoutePatterns := routePatternTemplates
	newRoutePatterns := make([]string, len(routePatternTemplates))

	// We must copy slice to avoid returning a reference to the original, underlying array on consecutive requests
	copy(newRoutePatterns, routePatternTemplates)

	for param, value := range paramMap {
		for i, pattern := range newRoutePatterns {
			newRoutePatterns[i] = strings.ReplaceAll(pattern, ":"+param, value)
		}
	}

	return newRoutePatterns
}

//* MEMORY
// Gets the allocated memory in bytes so we can compare it to the max allowed memory for the cache
// if that type of 'capacity' is chosen

// MemUsage outputs the current memory being used in bytes.
func MemUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	return m.Alloc
}

const (
	B uint64 = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)

var VALID_CAP_UNITS = []string{"B", "KB", "MB", "GB", "TB"}

// Converts kb, mb, gb, tb to bytes.
// If the unit is not set, it will return the size passed into the function as if they are already in bytes.
func ToBytes(size uint64, unit string) (uint64, error) {
	unit = strings.ToUpper(unit)

	switch unit {
	case "B":
		return size, nil
	case "KB":
		return size * KB, nil
	case "MB":
		return size * MB, nil
	case "GB":
		return size * GB, nil
	case "TB":
		return size * TB, nil
	default:
		err := fmt.Errorf("unknown capacity unit %q", unit)

		return 0, err
	}
}
