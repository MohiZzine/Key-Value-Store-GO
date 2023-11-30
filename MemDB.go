package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
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
	Empty     Error = iota
	threshold int   = 3
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
	wal                 *WAL
}

func NewMemDB() *MemDB {
	wal, err := NewWAL("wal")
	if err != nil {
		return nil
	}
	return &MemDB{
		sortedKeyValueStore: NewSortedKeyValueStore(),
		wal:                 wal,
	}
}

// Add a method to set the smallest and largest keys
func (mem *MemDB) setRangeKeys(smallestKey, largestKey string) {
	mem.smallestKey = smallestKey
	mem.largestKey = largestKey
}

// New function to check the threshold and flush data into SST files
func (mem *MemDB) checkAndFlush() error {
	fmt.Println("Checking threshold...")
	fmt.Println("Current key count:", len(mem.sortedKeyValueStore.keys))
	if len(mem.sortedKeyValueStore.keys) > threshold {
		// Increment the file index for naming
		mem.wal.Flush()

		// Flush the SortedKeyValueStore to an SST file
		keyValues := mem.sortedKeyValueStore.GetKeyValues()
		filename := "mohieddine_" + strconv.Itoa(mem.wal.currentIndex) + ".sst" // Adjust the naming convention as needed
		fmt.Println("Flushing to SST file:", filename)
		err := flushSSTFile(filename, keyValues)
		if err != nil {
			return err
		}

		// Clear the SortedKeyValueStore after flushing
		mem.sortedKeyValueStore = NewSortedKeyValueStore()
		mem.setRangeKeys("", "") // Reset range keys for the new SST file
	}
	return nil
}

func (mem *MemDB) Set(key, value string) error {
	mem.wal.WriteRecord(WALRecord{Operation: "Set", Key: key, Value: value})

	// Check if the key is within the range of keys in the SST file
	mem.sortedKeyValueStore.Set(key, value, true)

	// Check and flush if threshold is reached
	err := mem.checkAndFlush()
	if err != nil {
		return err
	}

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
	// Check if the key is within the range of keys in the SST files
	if (key < mem.smallestKey || key > mem.largestKey) && mem.smallestKey != "" && mem.largestKey != "" {
		return "", errors.New("Key probably in database")
	}

	// Retrieve the value and marker for the key from the SortedKeyValueStore
	valueMarkerPair, exists := mem.sortedKeyValueStore.Get(key)
	if exists == nil && valueMarkerPair != "" {
		// If the key is found in the SortedKeyValueStore and marked as present
		return valueMarkerPair, nil
	}

	// Check SST files from the most recent to the least recent
	for i := mem.wal.currentIndex; i >= 0; i-- {
		sstFile := "mohieddine_" + strconv.Itoa(i) + ".sst" // Adjust the naming convention as needed
		keyValues, _, _, err := parseSSTFile(sstFile)
		if err != nil {
			// Handle the error, possibly log it
			continue
		}

		// Iterate through key-values in the SST file
		for _, kv := range keyValues {
			if kv.Key == key {
				// If the key is found in the SST file, return the associated value
				return kv.Value, nil
			}
		}
	}

	// Key not found in MemDB or SST files
	return "", errors.New("Key not found")
}

func (mem *MemDB) Del(key string) (string, error) {
	mem.wal.WriteRecord(WALRecord{Operation: "Del", Key: key})

	// Check if the key is within the range of keys in the SST file
	val, err := mem.Get(key)
	if err != nil {
		return "", err
	}
	mem.sortedKeyValueStore.Set(key, val, false)
	return val, nil
}
