package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestWALWriteRecord(t *testing.T) {
	// Create a temporary WAL file for testing
	tmpfile, err := ioutil.TempFile("", "wal_test")
	if err != nil {
		t.Fatalf("Error creating temporary WAL file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	wal, err := NewWAL(tmpfile.Name())
	if err != nil {
		t.Fatalf("Error creating WAL: %v", err)
	}
	defer wal.Close()

	record := WALRecord{
		Operation: "Set",
		Key:       "testKey",
		Value:     "testValue",
		Timestamp: time.Now(),
	}

	err = wal.WriteRecord(record)
	if err != nil {
		t.Fatalf("Error writing record to WAL: %v", err)
	}
}

func TestWALFlush(t *testing.T) {
	// Create a temporary WAL file for testing
	tmpfile, err := ioutil.TempFile("", "wal_test")
	if err != nil {
		t.Fatalf("Error creating temporary WAL file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	wal, err := NewWAL(tmpfile.Name())
	if err != nil {
		t.Fatalf("Error creating WAL: %v", err)
	}
	defer wal.Close()

	// Write a record to the WAL
	record := WALRecord{
		Operation: "Set",
		Key:       "testKey",
		Value:     "testValue",
		Timestamp: time.Now(),
	}
	err = wal.WriteRecord(record)
	if err != nil {
		t.Fatalf("Error writing record to WAL: %v", err)
	}

	// Flush the WAL
	wal.Flush()

	// Check if the current index is incremented
	if wal.currentIndex != 1 {
		t.Errorf("Expected current index to be 1 after flush, got %d", wal.currentIndex)
	}
}
