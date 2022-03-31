package cache

import (
	"reflect"
	"testing"
)

var testData CacheData = CacheData{
	Headers: map[string]string{
		"Content-Type": "application/json",
		"X-Test":       "hi-mom",
	},
	Body: []byte(`{"test": "hi mom"}`),
}

func TestRequiredCapacity(t *testing.T) {
	defer func() { recover() }()

	New(0)

	t.Error("Expected cache.New to panic if the capacity is 0")
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

	if data == nil {
		t.Error("Expected cache.Get to return the entry set with cache.Set")
	}

	if !reflect.DeepEqual(*data, testData) {
		t.Error("Expected cache.Get to return the exact data that was set with cache.Set")
	}
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

	if data != nil {
		t.Errorf("Expected cache.Get to return nil when the entry does not exist, but returned: %+v", data)
	}
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

	size := cache.Size() // sanity check validates that this is also length of list, not just entries
	if size != 2 {
		t.Errorf("Expected cache size to not go above 2 when capacity is 2, but it is: %d", size)
	}

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

	matches := cache.Match([]string{"^/test[1-2]", "^/posts/.+"})
	// matches are not something on the cache, so we can't use sanityCheck with expectedKeys to validate
	if !reflect.DeepEqual(matches, []string{"/test1", "/test2", "/posts/4", "/posts/5"}) {
		t.Errorf("Expected cache.Match to return 4 matches, but it returned: %v", matches)
	}
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

	// Get all keys from list so we can compare them with map keys
	listKeys := make([]string, 0, len(cache.entries)) // init slice same capacity as map, since we expect those to match
	current := cache.lru
	for {
		if current == nil {
			break
		}

		listKeys = append(listKeys, current.key)
		current = current.next
	}

	// Lengths should be the same
	keysLen := len(cache.entries)
	listLen := len(listKeys)
	if keysLen != listLen {
		t.Errorf("Expected cache.entries to have same number of entries as the linked list, but the list has %d entries, and cache.entries has %d entries", listLen, keysLen)
	}

	// All keys should be the same, so we check from both map and list to see which are missing
	for _, key := range listKeys {
		if _, ok := cache.entries[key]; !ok {
			t.Errorf("Expected cache.entries to have entry for '%s', but it only exists in the linked list", key)
		}
	}

	for key := range cache.entries {
		if !elemInSlice(key, &listKeys) {
			t.Errorf("Expected linked list to have entry for '%s', but it only exists in cache.entries", key)
		}
	}
}

// Make sure that MRU and LRU prev / next pointers are set correctly.
// This includes cases when there are 0 entries, 1 entry, 2 entries and more than 2 entries.
func assertListEnds(t *testing.T, cache *LRUCache, expectedKeys []string) {
	t.Helper()

	listLen := listLength(cache)

	// 0 entries
	if listLen == 0 {
		if cache.lru != nil {
			t.Error("Expected cache.lru to be nil when there are no entries in cache, but it is not")
		}

		if cache.mru != nil {
			t.Error("Expected cache.mru to be nil when there are no entries in cache, but it is not")
		}

		return
	}

	// Make sure LRU and MRU are set to the first and last expected keys
	if cache.lru.key != expectedKeys[0] {
		t.Errorf("Expected cache.lru to be set to the first entry in the list, but it is set to '%s', where '%s' was expected", cache.lru.key, expectedKeys[0])
	}
	if cache.mru.key != expectedKeys[len(expectedKeys)-1] {
		t.Errorf("Expected cache.mru to be set to the last entry in the list, but it is set to '%s', where '%s' was expected", cache.mru.key, expectedKeys[len(expectedKeys)-1])
	}

	// LRU should always be the first entry
	if cache.lru.prev != nil {
		t.Errorf("Expected cache.lru.prev to always be nil, but it is: %+v", cache.lru.prev)
	}

	// MRU should always be the last entry
	if cache.mru.next != nil {
		t.Errorf("Expected cache.mru.next to always be nil, but it is: %+v", cache.mru.next)
	}

	// 1 entry
	if listLen == 1 {
		if cache.lru != cache.mru {
			t.Error("Expected cache.lru and cache.mru to be the same entry when there is only one entry, but they are not")
		}

		// lru.next can only be nil if there is only one entry
		if cache.lru.next != nil {
			t.Error("Expected cache.lru.next to be nil when it is the only entry, but it is not")
		}

		// mru.prev can only be nil if there is only one entry
		if cache.mru.prev != nil {
			t.Error("Expected cache.mru.prev to be nil when it is the only entry, but it is not")
		}

		return
	}

	// 2 entries
	if listLen == 2 {
		if cache.lru.next != cache.mru {
			t.Error("Expected cache.lru.next to be the same entry as cache.mru when they are the only two entries, but they are not")
		}

		if cache.mru.prev != cache.lru {
			t.Error("Expected cache.mru.prev to be the same entry as cache.lru when they are the only two entries, but they are not")
		}
		return
	}

	// lru and mru cannot point to each other unless there are only 2 entries
	if listLen > 2 {
		if cache.lru.next == cache.mru {
			t.Error("Expected cache.lru.next to not be the same entry as cache.mru when there are more than 2 entries, but they are the same entry")
		}

		if cache.mru.prev == cache.lru {
			t.Error("Expected cache.mru.prev to not be the same entry as cache.lru when there are more than 2 entries, but they are the same entry")
		}
		return
	}
}

// Also tests MRU and LRU are set correctly
func assertListOrder(t *testing.T, cache *LRUCache, expectedKeys []string) {
	t.Helper()

	current := cache.lru

	for entryNum, key := range expectedKeys {
		// Make sure that the expected keys array is not longer than the list
		if current == nil {
			t.Fatal("Expected list to have the same number of entries as the expected keys, but expectedKeys is longer")
			// return
		}

		// Make sure all entry keys are in the list
		if current.key != key {
			if entryNum == 0 {
				t.Errorf("Expected cache.lru to be '%s', but it is '%s'", key, current.key)
			} else if entryNum == len(expectedKeys)-1 {
				t.Errorf("Expected cache.mru to be '%s', but it is '%s'", key, current.key)
			} else {
				t.Errorf("Expected entry %d to be '%s', but it is '%s'", entryNum, key, current.key)
			}
		}
		current = current.next

		if entryNum == len(expectedKeys)-1 && current != nil {
			t.Fatal("Expected list to have as many entries as the expectedKeys, but list is longer")
			// return
		}
	}
}

// Make sure that cache entry map has all expected entries and nothing else
func assertKeys(t *testing.T, cache *LRUCache, expectedKeys []string) {
	t.Helper()

	if len(cache.entries) != len(expectedKeys) {
		t.Errorf("Expected cache.entries to have %d entries, but it has %d", len(expectedKeys), len(cache.entries))
	}

	for _, key := range expectedKeys {
		if _, ok := cache.entries[key]; !ok {
			t.Errorf("Expected cache.entries to have entry for '%s', but it does not", key)
		}
	}
}

func elemInSlice[T comparable](elem T, list *[]T) bool {
	for _, b := range *list {
		if b == elem {
			return true
		}
	}
	return false
}

func listLength(cache *LRUCache) int {
	length := 0
	for current := cache.lru; current != nil; current = current.next {
		length++
	}
	return length
}
