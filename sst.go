package main

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sort"
)

const (
	magicNumber uint64 = 0x6973656D
)

type SSTFile struct {
	file *os.File
}

func createSSTFile(filename string) (*SSTFile, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &SSTFile{file: file}, nil
}

func (s *SSTFile) close() error {
	return s.file.Close()
}

func parseSSTFile(filename string) ([]KeyValue, string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", "", err
	}
	defer file.Close()

	// Read and validate the magic number
	var magic uint64
	if err := binary.Read(file, binary.LittleEndian, &magic); err != nil {
		return nil, "", "", err
	}
	if magic != magicNumber {
		return nil, "", "", errors.New("Invalid SST file format")
	}

	// Read entry count
	var entryCount uint64
	if err := binary.Read(file, binary.LittleEndian, &entryCount); err != nil {
		return nil, "", "", err
	}

	// Read smallest key
	smallestKey, err := readString(file)
	if err != nil {
		return nil, "", "", err
	}

	// Read largest key
	largestKey, err := readString(file)
	if err != nil {
		return nil, "", "", err
	}

	// Read key-value pairs
	keyValues := make([]KeyValue, entryCount)
	for i := uint64(0); i < entryCount; i++ {
		key, err := readString(file)
		if err != nil {
			return nil, "", "", err
		}

		value, err := readString(file)
		if err != nil {
			return nil, "", "", err
		}

		keyValues[i] = KeyValue{Key: key, Value: value}
	}

	// TODO: Implement checksum validation

	return keyValues, smallestKey, largestKey, nil
}

func flushSSTFile(filename string, keyValues []KeyValue) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write magic number
	if err := binary.Write(file, binary.LittleEndian, magicNumber); err != nil {
		return err
	}

	// Write entry count
	entryCount := uint64(len(keyValues))
	if err := binary.Write(file, binary.LittleEndian, entryCount); err != nil {
		return err
	}

	// Sort keyValues by key
	sort.Slice(keyValues, func(i, j int) bool {
		return keyValues[i].Key < keyValues[j].Key
	})

	// Write smallest key
	if len(keyValues) > 0 {
		if err := writeString(file, keyValues[0].Key); err != nil {
			return err
		}
	}

	// Write largest key
	if len(keyValues) > 0 {
		if err := writeString(file, keyValues[len(keyValues)-1].Key); err != nil {
			return err
		}
	}

	// Write key-value pairs
	for _, kv := range keyValues {
		if err := writeString(file, kv.Key); err != nil {
			return err
		}
		if err := writeString(file, kv.Value); err != nil {
			return err
		}
	}

	// TODO: Calculate and write checksum

	return nil
}

func writeString(w io.Writer, s string) error {
	// Write string length
	if err := binary.Write(w, binary.LittleEndian, uint64(len(s))); err != nil {
		return err
	}
	// Write string data
	_, err := w.Write([]byte(s))
	return err
}

func readString(r io.Reader) (string, error) {
	// Read string length
	var length uint64
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	// Read string data
	data := make([]byte, length)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
