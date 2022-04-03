package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBytes(t *testing.T) {
	assert := assert.New(t)

	type toBytesArgs struct {
		size     uint64
		unit     string
		expected uint64
	}

	testCases := [...]toBytesArgs{
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

	for _, c := range testCases {
		bytes, err := ToBytes(c.size, c.unit)

		if c.unit == "nb" || c.unit == "NB" || c.unit == "" {
			assert.Error(err, "Expected unknown unit \"%s\" to return an error", c.unit)
			assert.Zero(bytes)
		} else {
			assert.NoError(err)
			assert.Equal(c.expected, bytes, "Expected ToBytes(%d, \"%s\") to return %d, it returned %d", c.size, c.unit, c.expected, bytes)
		}
	}
}

func TestSetInitEmpty(t *testing.T) {
	set := make(Set[string])

	assert.Empty(t, set, "Expected set to be empty after initialization")
}

func TestSetAddAndHas(t *testing.T) {
	assert := assert.New(t)

	set := make(Set[string])

	set.Add("a")

	assert.Len(set, 1, "Expected set to increase length to 1 after adding an element, but length is now %d", len(set))

	assert.True(set.Has("a"), "Expected set.Has to return true when the set includes the element: \"a\"")
	assert.Contains(set, "a", "Expected set's underlying map to contain the element: \"a\"")

	assert.False(set.Has("b"), "Expected set.Has to return false when the set does not include the element: \"b\"")
	assert.NotContains(set, "b", "Expected set's underlying map to not contain the element: \"b\"")
}

func TestSetRemove(t *testing.T) {
	set := make(Set[string])

	set.Add("a")
	set.Remove("a")

	assert.False(t, set.Has("a"), "Expected set.Remove() to remove the element: \"a\"")
	assert.NotContains(t, "a", "Expected set's underlying map to not contain the element: \"a\"")
}

func TestSetElements(t *testing.T) {
	set := make(Set[string])

	set.Add("a")
	set.Add("b")

	assert.ElementsMatch(t, set.Elements(), []string{"a", "b"}, "Expected set.Elements to return the elements: [ \"a\", \"b\" ]")
	assert.Len(t, set, 2, "Expected set to only include 2 elements")
}
