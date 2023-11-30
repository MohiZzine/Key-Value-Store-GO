package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// WAL represents the Write-Ahead Log.
type WAL struct {
	file         *os.File
	mu           sync.Mutex
	currentIndex int // New field to track the current index
	watermark    int // New field to track the last successfully flushed index
}

// NewWAL creates a new Write-Ahead Log.
func NewWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file:         file,
		currentIndex: 0, // Initialize the current index to 0
		watermark:    0, // Initialize the watermark to 0
	}, nil
}

// WriteRecord writes a WALRecord to the Write-Ahead Log.
func (wal *WAL) WriteRecord(record WALRecord) error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Serialize the record to JSON
	jsonRecord, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Append the JSON record to the WAL file
	_, err = fmt.Fprintln(wal.file, string(jsonRecord))
	if err != nil {
		return err
	}

	return nil
}

// Flush resets the current index and updates the watermark.
func (wal *WAL) Flush() {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Update the watermark to the current index
	wal.watermark = wal.currentIndex

	// Reset the current index
	wal.currentIndex++
}

// Close closes the Write-Ahead Log file.
func (wal *WAL) Close() error {
	return wal.file.Close()
}
