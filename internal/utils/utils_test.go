package utils

import (
	"testing"
)

func TestToBytes(t *testing.T) {
	type toBytesArgs struct {
		size     uint64
		unit     string
		expected uint64
	}

	testCases := [...]toBytesArgs{
		{size: 1, unit: "kb", expected: 1024},
		{size: 4, unit: "KB", expected: 1024 * 4},
		{size: 2, unit: "mb", expected: 2097152},
		{size: 1, unit: "MB", expected: 1048576},
		{size: 1, unit: "MB", expected: 1048576},
		{size: 1, unit: "gb", expected: 1073741824},
		{size: 5, unit: "GB", expected: 5368709120},
		{size: 1, unit: "tb", expected: 1099511627776},
		{size: 2, unit: "TB", expected: 2199023255552},
		{size: 1, unit: "nb", expected: 1},
		{size: 999, unit: "NB", expected: 999},
		{size: 666, unit: "", expected: 666},
	}

	for _, c := range testCases {
		bytes := ToBytes(c.size, c.unit)

		if bytes != c.expected {
			t.Errorf("Expected ToBytes(%d, %s) to return %d, it returned %d", c.size, c.unit, c.expected, bytes)
		}
	}
}

func TestSetInitEmpty(t *testing.T) {
	set := make(Set[string])

	if len(set) != 0 {
		t.Errorf("Expected Set to be initialized empty, but it has %d elements", len(set))
	}
}

func TestSetAddAndHas(t *testing.T) {
	set := make(Set[string])

	set.Add("a")
	if len(set) != 1 {
		t.Errorf("Expected set to increase length to 1 after adding an element, but length is now %d", len(set))
	}
	if set.Has("a") == false {
		t.Errorf("Expected set.Has() to return true when the set includes the given element")
	}
	if set.Has("b") == true {
		t.Errorf("Expected set.Has() to return false when the set does not include the given element")
	}
}

func TestSetRemove(t *testing.T) {
	set := make(Set[string])

	set.Add("a")
	set.Remove("a")
	if set.Has("a") == true {
		t.Errorf("Expected set.Remove() to remove the given element")
	}
}

func TestSetElements(t *testing.T) {
	set := make(Set[string])

	set.Add("a")
	set.Add("b")
	if sliceEqual(t, set.Elements(), []string{"a", "b"}) == false {
		t.Errorf("Expected set.Elements() to return the elements in the set in the order they were added, but returned %v", set.Elements())
	}
}

// sliceEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func sliceEqual[T comparable](t *testing.T, a, b []T) bool {
	t.Helper()

	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
