package main

// import (
// 	"encoding/binary"
// )

// func (mem *memDB) writeToLogFile(cmd Cmd, key, value string) error {
// 	// Binary encode the log entry (marker + key + value)
// 	entry := make([]byte, binary.MaxVarintLen64*3+len(key)+len(value))
// 	pos := binary.PutVarint(entry, int64(cmd))
// 	pos += binary.PutVarint(entry[pos:], int64(len(key)))
// 	copy(entry[pos:], key)
// 	pos += len(key)
// 	binary.PutVarint(entry[pos:], int64(len(value)))
// 	copy(entry[pos+len(value):], value)

// 	_, err := mem.logFile.Write(entry)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // CloseLogFile closes the log file when the program exits
// func (mem *memDB) CloseLogFile() {
// 	mem.logFile.Close()
// }

// func main() {
// 	// Instantiate memDB with a log file
// 	memDB, err := NewInMem("log.bin")
// 	if err != nil {
// 		fmt.Println("Error creating memDB:", err)
// 		return
// 	}
// 	defer memDB.CloseLogFile()

// 	// Use memDB as needed (set, get, del)
// 	memDB.Set("key1", "value1")
// 	memDB.Set("key2", "value2")
// 	val, err := memDB.Get("key1")
// 	if err != nil {
// 		fmt.Println("Error getting key1:", err)
// 	} else {
// 		fmt.Println("Get result:", val)
// 	}

// 	memDB.Del("key1")

// 	val, err = memDB.Get("key1")
// 	if err != nil {
// 		fmt.Println("Error getting key1:", err)
// 	} else {
// 		fmt.Println("Get result:", val)
// 	}
// }
