package utils

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/NormalReedus/cache-me-ousside/internal/logger"
)

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
func ToBytes(size uint64, unit string) (uint64, error) {
	unit = strings.ToLower(unit)

	switch unit {
	case "b":
		return size, nil
	case "kb":
		return size << 10, nil
	case "mb":
		return size << 20, nil
	case "gb":
		return size << 30, nil
	case "tb":
		return size << 40, nil
	default:
		err := fmt.Errorf("unknown unit: %q", unit)

		logger.Error(err)

		return 0, err
	}
}
