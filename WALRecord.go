// WALRecord.go

package main

import (
	"encoding/json"
	"time"
)

// Constants for operations
const (
	SetOperation = "Set"
	DelOperation = "Del"
)

// WALRecord represents a record in the Write-Ahead Log.
type WALRecord struct {
	Operation string    `json:"operation"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// NewWALRecord creates a new WALRecord.
func NewWALRecord(operation, key, value string) WALRecord {
	return WALRecord{
		Operation: operation,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}
}

// NewSetWALRecord creates a new WALRecord for a 'Set' operation.
func NewSetWALRecord(key, value string) WALRecord {
	return NewWALRecord(SetOperation, key, value)
}

// NewDelWALRecord creates a new WALRecord for a 'Del' operation.
func NewDelWALRecord(key string) WALRecord {
	return NewWALRecord(DelOperation, key, "")
}

// Serialize serializes the WALRecord to JSON.
func (r *WALRecord) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

// Deserialize deserializes JSON data into a WALRecord.
func Deserialize(data []byte) (*WALRecord, error) {
	var record WALRecord
	err := json.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
