package cache

import (
	"encoding/json"
	"fmt"
	"log"
)

func newEntry(key string, data *CacheData) *Entry {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(fmt.Errorf("there was an creating the entry: %v", key))
		return nil
	}

	entry := &Entry{
		key:  key,
		data: jsonData,
	}

	return entry
}

type Entry struct {
	key  string
	data []byte // marshaled json of CacheData
	next *Entry
	prev *Entry
}

//! TEST
func (entry *Entry) Key() string {
	return entry.key
}

//! TESTEND

func (entry *Entry) setNext(newEntry *Entry) *Entry {
	if newEntry == nil {
		return nil
	}

	// If this entry is not head, insert newEntry between this entry and next entry
	nextEntry := entry.next
	entry.next = newEntry
	newEntry.prev = entry
	// if entry is head, nextEntry is nil and will correctly set the newEntry's next as nil, because newEntry should become head
	newEntry.next = nextEntry
	// If nextEntry is not nil (entry is not head), set the prev of nextEntry to newEntry
	if nextEntry != nil {
		nextEntry.prev = newEntry
	}

	return newEntry
}

func (entry *Entry) unmarshalData() CacheData {
	var data CacheData
	json.Unmarshal(entry.data, &data)

	return data
}
