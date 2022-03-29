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

	cases := [...]toBytesArgs{
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

	for _, c := range cases {
		bytes := ToBytes(c.size, c.unit)
		if bytes != c.expected {
			t.Errorf("ToBytes(%d, %s) returned %d, expected %d", c.size, c.unit, bytes, c.expected)
		}
	}
}

func TestSet(t *testing.T) {
	set := make(Set[string])

	if len(set) != 0 {
		t.Errorf("Set should be initialized empty, but has %d elements", len(set))
	}

	set.Add("a")
	if len(set) != 1 {
		t.Errorf("After adding an element, set should increase length to 1, but length is now %d", len(set))
	}
	if set.Has("a") == false {
		t.Errorf("set.Has() should return true when the set includes the given element")
	}
	if set.Has("b") == true {
		t.Errorf("set.Has() should return false when the set does not include the given element")
	}

	set.Add("b")
	set.Remove("a")
	if set.Has("a") == true {
		t.Errorf("set.Remove() should remove the given element")
	}

	set.Add("c")
	if sliceEqual(set.Elements(), []string{"c", "b"}) == false {
		t.Errorf("set.Elements() should return the elements in the set in the order they were added, but returned %v", set.Elements())
	}
}

// sliceEqual tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func sliceEqual[T comparable](a, b []T) bool {
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
