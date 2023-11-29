package main

import (
	"errors"
	"sort"
)

// KeyValue represents a key-value pair.
type KeyValue struct {
	Key   string
	Value string
}

type Cmd int

const (
	Get Cmd = iota
	Set
	Del
	Ext
	Unk
)

type Error int

func (e Error) Error() string {
	return "Empty command"
}

const (
	Empty Error = iota
)

type DB interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Del(key string) (string, error)
}

type ValueMarkerPair struct {
	Value  string
	Marker bool
}

type SortedKeyValueStore struct {
	values  map[string]ValueMarkerPair
	keys    []string
	markers map[string]bool
}

func NewSortedKeyValueStore() *SortedKeyValueStore {
	return &SortedKeyValueStore{
		values:  make(map[string]ValueMarkerPair),
		keys:    make([]string, 0),
		markers: make(map[string]bool),
	}
}

func (store *SortedKeyValueStore) Set(key string, value string, marker bool) {
	// If the key already exists, update the value and marker
	if _, exists := store.values[key]; exists {
		store.values[key] = ValueMarkerPair{Value: value, Marker: marker}
	} else {
		// Otherwise, add the new key
		store.keys = append(store.keys, key)
		store.values[key] = ValueMarkerPair{Value: value, Marker: marker}
		store.markers[key] = marker
		sort.Strings(store.keys)
	}
}

func (store *SortedKeyValueStore) Get(key string) (string, error) {
	// Check if the key exists
	if _, exists := store.values[key]; !exists {
		return "", errors.New("Key probably in database")
	}

	// Retrieve the value and marker for the key
	valueMarkerPair := store.values[key]
	if valueMarkerPair.Marker {
		// If marker is true, return the value
		return valueMarkerPair.Value, nil
	} else {
		// If marker is false, return "key not found" error
		return "", errors.New("Key not found")
	}
}

// Load loads key-values into the SortedKeyValueStore.
func (store *SortedKeyValueStore) Load(keyValues []KeyValue) {
	for _, kv := range keyValues {
		store.Set(kv.Key, kv.Value, true)
	}
}

func (store *SortedKeyValueStore) GetKeyValues() []KeyValue {
	keyValues := make([]KeyValue, 0, len(store.keys))

	for _, key := range store.keys {
		valueMarkerPair := store.values[key]
		keyValues = append(keyValues, KeyValue{Key: key, Value: valueMarkerPair.Value})
	}

	return keyValues
}

type MemDB struct {
	sortedKeyValueStore *SortedKeyValueStore
	smallestKey         string
	largestKey          string
}

func NewMemDB() *MemDB {
	return &MemDB{
		sortedKeyValueStore: NewSortedKeyValueStore(),
	}
}

// Add a method to set the smallest and largest keys
func (mem *MemDB) setRangeKeys(smallestKey, largestKey string) {
	mem.smallestKey = smallestKey
	mem.largestKey = largestKey
}

func (mem *MemDB) Set(key, value string) error {
	mem.sortedKeyValueStore.Set(key, value, true)
	return nil
}

func (mem *MemDB) LoadSSTFile(filename string) error {
	keyValues, smallestKey, largestKey, err := parseSSTFile(filename)
	if err != nil {
		return err
	}

	mem.sortedKeyValueStore.Load(keyValues)
	mem.setRangeKeys(smallestKey, largestKey)
	return nil
}

func (mem *MemDB) Get(key string) (string, error) {
	// Check if the key is within the range of keys in the SST file
	if key < mem.smallestKey || key > mem.largestKey {
		return "", errors.New("Key probably in database")
	}

	// Retrieve the value and marker for the key from the SortedKeyValueStore
	valueMarkerPair, exists := mem.sortedKeyValueStore.Get(key)
	if exists != nil {
		return "", errors.New("Key not found")
	}

	if valueMarkerPair != "" {
		// If marker is true, return the value
		return valueMarkerPair, nil
	} else {
		// If marker is false, return "key not found" error
		return "", errors.New("Key has been deleted")
	}
}

func (mem *MemDB) Del(key string) (string, error) {
	val, err := mem.Get(key)
	if err != nil {
		return "", err
	}
	mem.sortedKeyValueStore.Set(key, val, false)
	return val, nil
}
