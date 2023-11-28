package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// MemDB represents an in-memory database.
type MemDB struct {
	mu     sync.RWMutex
	values map[string]string
}

// Set sets the value for the given key in the MemDB.
func (mem *MemDB) Set(key, value string) {
	mem.mu.Lock()
	defer mem.mu.Unlock()
	mem.values[key] = value
}

// Get gets the value for the given key from the MemDB.
func (mem *MemDB) Get(key string) (string, bool) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()
	val, ok := mem.values[key]
	return val, ok
}

// Del deletes the value for the given key from the MemDB.
func (mem *MemDB) Del(key string) (string, bool) {
	mem.mu.Lock()
	defer mem.mu.Unlock()
	val, ok := mem.values[key]
	if ok {
		delete(mem.values, key)
	}
	return val, ok
}

// NewMemDB creates a new instance of MemDB.
func NewMemDB() *MemDB {
	return &MemDB{
		values: make(map[string]string),
	}
}

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
	val, ok := l.MemDB.Get(key)
	if !ok {
		return "Value Probably in the database", nil
	}
	return val, nil
}

// Del deletes the value for the given key from the LSTM's MemDB.
func (l *LSTM) Del(key string) (string, error) {
	val, ok := l.MemDB.Del(key)
	if !ok {
		return "", fmt.Errorf("Key not found")
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
	memDB := NewMemDB()
	lstm := &LSTM{MemDB: memDB}
	handler := &Handler{db: lstm}

	http.HandleFunc("/get", handler.GetHandler)
	http.HandleFunc("/set", handler.SetHandler)
	http.HandleFunc("/del", handler.DelHandler)

	fmt.Println("Server is listening on :8080")
	http.ListenAndServe(":8080", nil)
}
