package main

import (
	"strconv"
	"testing"
)

func TestMemDBSetGet(t *testing.T) {
	memDB := NewMemDB()

	key := "testKey"
	value := "testValue"

	err := memDB.Set(key, value)
	if err != nil {
		t.Fatalf("Error setting key-value pair: %v", err)
	}

	result, err := memDB.Get(key)
	if err != nil {
		t.Fatalf("Error getting value for key: %v", err)
	}

	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}
}

func TestMemDBDel(t *testing.T) {
	memDB := NewMemDB()

	key := "testKey"
	value := "testValue"

	err := memDB.Set(key, value)
	if err != nil {
		t.Fatalf("Error setting key-value pair: %v", err)
	}

	result, err := memDB.Del(key)
	if err != nil {
		t.Fatalf("Error deleting key-value pair: %v", err)
	}

	if result != value {
		t.Errorf("Expected deleted value %s, got %s", value, result)
	}

	_, err = memDB.Get(key)
	if err == nil {
		t.Error("Expected error for Get after deletion, but got nil")
	}
}

func TestMemDBThresholdFlush(t *testing.T) {
	memDB := NewMemDB()
	threshold := 3

	for i := 1; i <= threshold+1; i++ {
		key := strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		err := memDB.Set(key, value)
		if err != nil {
			t.Fatalf("Error setting key-value pair: %v", err)
		}
	}

	keyValues := memDB.sortedKeyValueStore.GetKeyValues()
	if len(keyValues) != 0 {
		t.Errorf("Expected SortedKeyValueStore to be empty after threshold flush, got %d items", len(keyValues))
	}
}
