package main

import (
	"errors"
	"sort"
)

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

type MemDB struct {
	sortedKeyValueStore *SortedKeyValueStore
}

func NewMemDB() *MemDB {
	return &MemDB{
		sortedKeyValueStore: NewSortedKeyValueStore(),
	}
}

func (mem *MemDB) Set(key, value string) error {
	mem.sortedKeyValueStore.Set(key, value, true)
	return nil
}

func (mem *MemDB) Get(key string) (string, error) {
	val, err := mem.sortedKeyValueStore.Get(key)
	if err != nil {

		return "", err
	}

	return val, nil
}

func (mem *MemDB) Del(key string) (string, error) {
	val, err := mem.Get(key)
	if err != nil {
		return "", err
	}
	mem.sortedKeyValueStore.Set(key, val, false)
	return val, nil
}
