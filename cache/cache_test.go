package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = CacheData{
	Headers: map[string]string{
		"Content-Type": "application/json",
		"X-Test":       "hi-mom",
	},
	Body: []byte(`{"test": "hi mom"}`),
}

func TestRequiredCapacity(t *testing.T) {
	defer func() { recover() }()

	assert.Panics(t, func() { New(0) }, "Expected cache.New to panic if the capacity is 0")
}

func TestSetEntry(t *testing.T) {
	cache := New(1)

	cache.Set("/test1", &testData)

	expectedKeys := []string{"/test1"}

	sanityCheck(t, cache, expectedKeys)
}

func TestEntriesOrder(t *testing.T) {
	cache := New(5)

	cache.Set("/test1", &testData)
	cache.Set("/test2", &testData)
	cache.Set("/test3", &testData)
	cache.Set("/posts/4", &testData)
	cache.Set("/posts/5", &testData)

	expectedKeys := []string{
		"/test1",
		"/test2",
		"/test3",
		"/posts/4",
		"/posts/5",
	}
	sanityCheck(t, cache, expectedKeys) // implicitly also checks that Set will add things as MRU

	// /test1 should now be moved to MRU
	cache.Get("/test1")

	expectedKeys = []string{
		"/test2",
		"/test3",
		"/posts/4",
		"/posts/5",
		"/test1",
	}
	sanityCheck(t, cache, expectedKeys)
}

func TestGetEntry(t *testing.T) {
	cache := New(1)

	cache.Set("/test1", &testData)

	data := cache.Get("/test1")

	expectedKeys := []string{
		"/test1",
	}
	sanityCheck(t, cache, expectedKeys)

	assert.NotNil(t, data, "Expected cache.Get to return the entry set with cache.Set")

	assert.Equal(t, &testData, data, "Expected cache.Get to return the exact data that was set with cache.Set")
}

func TestGetMissingEntry(t *testing.T) {
	cache := New(1)

	cache.Set("/test1", &testData)

	data := cache.Get("/test2")

	// Sanity check before other tests
	expectedKeys := []string{
		"/test1",
	}
	sanityCheck(t, cache, expectedKeys)

	assert.Nil(t, data, "Expected cache.Get to return nil when the entry does not exist, but returned: %+v", data)
}

func TestEvict(t *testing.T) {
	// Cannot really test for memory based eviction since memory usage will vary
	cache := New(2)

	cache.Set("/test1", &testData)
	cache.Set("/test2", &testData)
	cache.Set("/test3", &testData)

	// Sanity check before other tests
	expectedKeys := []string{
		"/test2",
		"/test3",
	}
	sanityCheck(t, cache, expectedKeys)

	size := cache.Size() // sanity check validates that Size() is also length of list, not just entries

	assert.Equal(t, 2, size, "Expected cache size to not go above 2 when capacity is 2, but it is: %d", size)

	expectedKeys = []string{
		"/test2",
		"/test3",
	}
	sanityCheck(t, cache, expectedKeys)
}

func TestMatch(t *testing.T) {
	cache := New(5)

	cache.Set("/test1", &testData)
	cache.Set("/test2", &testData)
	cache.Set("/test3", &testData)
	cache.Set("/posts/4", &testData)
	cache.Set("/posts/5", &testData)

	// Sanity check before other tests
	expectedKeys := []string{
		"/test1",
		"/test2",
		"/test3",
		"/posts/4",
		"/posts/5",
	}
	sanityCheck(t, cache, expectedKeys)

	patterns := []string{
		"^/test[1-2]",
		"^/posts/.+",
	}
	expectedMatches := []string{
		"/test1",
		"/test2",
		"/posts/4",
		"/posts/5",
	}

	matches := cache.Match(patterns)
	// matches are not something on the cache, so we can't use sanityCheck with expectedKeys to validate
	assert.ElementsMatch(t, expectedMatches, matches, "Expected cache.Match to return the the keys %v, but returned %v", expectedMatches, matches)
}

func TestBust(t *testing.T) {
	cache := New(5)

	cache.Set("/test1", &testData)
	cache.Set("/test2", &testData)
	cache.Set("/test3", &testData)
	cache.Set("/posts/4", &testData)
	cache.Set("/posts/5", &testData)

	// Sanity check before other tests
	expectedKeys := []string{
		"/test1",
		"/test2",
		"/test3",
		"/posts/4",
		"/posts/5",
	}
	sanityCheck(t, cache, expectedKeys)

	// Bust over two passes just to try out both single and multi busting
	cache.Bust("/test1", "/test2")
	cache.Bust("/posts/5")

	expectedKeys = []string{
		"/test3",
		"/posts/4",
	}
	sanityCheck(t, cache, expectedKeys)
}

//* USEFUL ASSERTIONS

// Bundles all the following utility assertions, so we can easily always check the cache state is alright
func sanityCheck(t *testing.T, cache *LRUCache, expectedKeys []string) {
	t.Helper()

	assertEntriesMatch(t, cache)
	assertListEnds(t, cache, expectedKeys)
	assertListOrder(t, cache, expectedKeys)
	assertKeys(t, cache, expectedKeys)
}

// Make sure that they same keys are saved in the list as is in the cache map, and nothing else.
// Also makes sure lengths of both map and list are the same.
func assertEntriesMatch(t *testing.T, cache *LRUCache) {
	t.Helper()
	assert := assert.New(t)

	// Get all keys from list so we can compare them with map keys
	listKeys := listKeys(cache)

	// Lengths should be the same
	keysLen := len(cache.entries)
	listLen := len(listKeys)

	assert.Equal(keysLen, listLen, "Expected cache.entries to have same number of entries as the linked list, but the list has %d entries, and cache.entries has %d entries", listLen, keysLen)

	// All keys should be the same, so we check from both map and list to see which are missing
	assert.ElementsMatch(listKeys, cache.CachedEndpoints(), "Expected cache.entries to have the same keys as the linked list, but the list has %v, and cache.entries has %v", listKeys, cache.CachedEndpoints())
}

// Make sure that MRU and LRU prev / next pointers are set correctly.
// This includes cases when there are 0 entries, 1 entry, 2 entries and more than 2 entries.
func assertListEnds(t *testing.T, cache *LRUCache, expectedKeys []string) {
	t.Helper()
	assert := assert.New(t)

	listLen := listLength(cache)

	// 0 entries
	if listLen == 0 {
		assert.Nil(cache.lru, "Expected cache.lru to be nil when there are no entries in cache, but it is not")

		assert.Nil(cache.mru, "Expected cache.mru to be nil when there are no entries in cache, but it is not")

		return
	}

	// Make sure LRU and MRU are set to the first and last expected keys
	assert.Equal(expectedKeys[0], cache.lru.key, "Expected cache.lru to be set to the first entry in the list, but it is set to '%s', where '%s' was expected", cache.lru.key, expectedKeys[0])

	assert.Equal(expectedKeys[len(expectedKeys)-1], cache.mru.key, "Expected cache.mru to be set to the last entry in the list, but it is set to '%s', where '%s' was expected", cache.mru.key, expectedKeys[len(expectedKeys)-1])

	// LRU should always be the first entry
	assert.Nil(cache.lru.prev, "Expected cache.lru.prev to always be nil, but it is: %+v", cache.lru.prev)

	// MRU should always be the last entry
	assert.Nil(cache.mru.next, "Expected cache.mru.next to always be nil, but it is: %+v", cache.mru.next)

	// 1 entry
	if listLen == 1 {
		assert.Same(cache.lru, cache.mru, "Expected cache.lru and cache.mru to be the same entry when there is only one entry, but they are not")

		// lru.next can only be nil if there is only one entry
		assert.Nil(cache.lru.next, "Expected cache.lru.next to be nil when it is the only entry, but it is: %+v", cache.lru.next)

		// mru.prev can only be nil if there is only one entry
		assert.Nil(cache.mru.next, "Expected cache.mru.next to be nil when it is the only entry, but it is: %+v", cache.mru.next)

		return
	}

	// 2 entries
	if listLen == 2 {
		assert.Same(cache.lru.next, cache.mru, "Expected cache.lru.next to be the same entry as cache.mru when they are the only two entries, but they are not")

		assert.Same(cache.mru.prev, cache.lru, "Expected cache.mru.prev to be the same entry as cache.lru when they are the only two entries, but they are not")

		return
	}

	// lru and mru cannot point to each other unless there are only 2 entries
	if listLen > 2 {
		assert.NotSame(cache.lru.next, cache.mru, "Expected cache.lru.next to not be the same entry as cache.mru when there are more than 2 entries, but they are the same entry")

		assert.NotSame(cache.mru.prev, cache.lru, "Expected cache.mru.prev to not be the same entry as cache.lru when there are more than 2 entries, but they are the same entry")

		return
	}
}

// Also tests MRU and LRU are set correctly
func assertListOrder(t *testing.T, cache *LRUCache, expectedKeys []string) {
	t.Helper()
	assert := assert.New(t)

	// Get all keys from list so we can compare them with expected order
	listKeys := listKeys(cache)

	assert.Equal(len(expectedKeys), len(listKeys), "Expected list to have the same number of entries as the expected keys, got %d, but expected %d", len(listKeys), len(expectedKeys))

	assert.Equal(expectedKeys, listKeys, "Expected list to have the same keys as the expected keys, got %v, but expected %v", listKeys, expectedKeys)
}

// Make sure that cache entry map has all expected entries and nothing else
func assertKeys(t *testing.T, cache *LRUCache, expectedKeys []string) {
	t.Helper()
	assert := assert.New(t)

	assert.Equal(len(expectedKeys), len(cache.entries), "Expected cache.entries to have %d entries, but it has %d", len(expectedKeys), len(cache.entries))

	assert.ElementsMatch(expectedKeys, cache.CachedEndpoints(), "Expected cache to have entries for %v, but got %v", expectedKeys, cache.CachedEndpoints())
}

func listLength(cache *LRUCache) int {
	length := 0
	for current := cache.lru; current != nil; current = current.next {
		length++
	}
	return length
}

func listKeys(cache *LRUCache) []string {
	keys := make([]string, 0, len(cache.entries))

	for current := cache.lru; current != nil; current = current.next {
		keys = append(keys, current.key)
	}

	return keys
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
