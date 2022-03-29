package utils

import (
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

//* MEMORY
// Gets the allocated memory in bytes so we can compare it to the max allowed memory for the cache
// if that type of 'capacity' is chosen

// MemUsage outputs the currentmemory being used in bytes.
func MemUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	return m.Alloc
}

// Converts kb, mb, gb, tb to bytes.
// If the unit is not set, it will return the size passed into the function as if they are already in bytes.
func ToBytes(size uint64, unit string) uint64 {
	unit = strings.ToLower(unit)

	switch unit {
	case "kb":
		return size << 10
	case "mb":
		return size << 20
	case "gb":
		return size << 30
	case "tb":
		return size << 40
	default:
		return size
	}
}
