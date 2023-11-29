package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// WALRecord represents a record in the Write-Ahead Log.
type WALRecord struct {
	Operation string    `json:"operation"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// WAL represents the Write-Ahead Log.
type WAL struct {
	file *os.File
	mu   sync.Mutex
}

// NewWAL creates a new Write-Ahead Log.
func NewWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file: file,
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

// Close closes the Write-Ahead Log file.
func (wal *WAL) Close() error {
	return wal.file.Close()
}
