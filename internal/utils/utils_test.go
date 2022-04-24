package utils

import (
	"testing"

	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.Initialize("")
}

func TestToBytes(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		size     uint64
		unit     string
		expected uint64
	}

	tests := [...]args{
		{size: 12345, unit: "b", expected: 12345},
		{size: 6789, unit: "B", expected: 6789},
		{size: 1, unit: "kb", expected: 1024},
		{size: 4, unit: "KB", expected: 1024 * 4},
		{size: 2, unit: "mb", expected: 2097152},
		{size: 1, unit: "MB", expected: 1048576},
		{size: 1, unit: "MB", expected: 1048576},
		{size: 1, unit: "gb", expected: 1073741824},
		{size: 5, unit: "GB", expected: 5368709120},
		{size: 1, unit: "tb", expected: 1099511627776},
		{size: 2, unit: "TB", expected: 2199023255552},
		{size: 1, unit: "nb", expected: 0},
		{size: 999, unit: "NB", expected: 0},
		{size: 666, unit: "", expected: 0},
	}

	for _, tt := range tests {
		bytes, err := ToBytes(tt.size, tt.unit)

		if tt.unit == "nb" || tt.unit == "NB" || tt.unit == "" {
			assert.Error(err, "Expected unknown unit %q to return an error", tt.unit)
			assert.Zero(bytes)
		} else {
			assert.NoError(err)
			assert.Equal(tt.expected, bytes, "Expected ToBytes(%d, %q) to return %d, it returned %d", tt.size, tt.unit, tt.expected, bytes)
		}
	}
}
