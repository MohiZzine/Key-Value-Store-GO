package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// LSTM represents a key-value store that uses an in-memory database.
type LSTM struct {
	MemDB *MemDB
}

// Set sets the value for the given key in the LSTM's MemDB.
func (l *LSTM) Set(key, value string) error {
	l.MemDB.Set(key, value)
	return nil
}

// Get gets the value for the given key from the LSTM's MemDB.
// If the key is not found, it returns "Value Probably in the database".
func (l *LSTM) Get(key string) (string, error) {
	val, err := l.MemDB.Get(key)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Del deletes the value for the given key from the LSTM's MemDB.
func (l *LSTM) Del(key string) (string, error) {
	val, err := l.MemDB.Del(key)
	if err != nil {
		return "", err
	}
	return val, nil
}

type Handler struct {
	db DB
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, err := h.db.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, value)
}

func (h *Handler) SetHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, ok := data["key"]
	if !ok {
		http.Error(w, "Key is required in the JSON body", http.StatusBadRequest)
		return
	}

	value, ok := data["value"]
	if !ok {
		http.Error(w, "Value is required in the JSON body", http.StatusBadRequest)
		return
	}

	if err := h.db.Set(key, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DelHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, err := h.db.Del(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, value)
}

func main() {
	// memDB := NewMemDB()
	// lstm := &LSTM{MemDB: memDB}
	// handler := &Handler{db: lstm}

	// http.HandleFunc("/get", handler.GetHandler)
	// http.HandleFunc("/set", handler.SetHandler)
	// http.HandleFunc("/del", handler.DelHandler)

	// fmt.Println("Server is listening on :8080")
	// http.ListenAndServe(":8080", nil)

	// Create a new MemDB
	memDB := NewMemDB()
	lstm := &LSTM{MemDB: memDB}
	handler := &Handler{db: lstm}

	// Define routes
	http.HandleFunc("/get", handler.GetHandler)
	http.HandleFunc("/set", handler.SetHandler)
	http.HandleFunc("/del", handler.DelHandler)

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server is listening on :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Perform some operations on the database
	memDB.Set("key1", "value1")
	memDB.Set("key2", "value2")
	memDB.Set("key3", "value3")

	// Flush the data to an SST file
	err := flushSSTFile("data.sst", memDB.sortedKeyValueStore.GetKeyValues())
	if err != nil {
		fmt.Println("Error flushing data to SST file:", err)
		os.Exit(1)
	}

	// Close the server after performing operations
	fmt.Println("Press Ctrl+C to stop the server...")
	select {}
}
