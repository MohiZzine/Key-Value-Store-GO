package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPISetGet(t *testing.T) {
	lstm := &LSTM{MemDB: NewMemDB()}
	handler := &Handler{db: lstm}

	key := "testKey"
	value := "testValue"

	setRequest := map[string]string{"key": key, "value": value}
	setBody, err := json.Marshal(setRequest)
	if err != nil {
		t.Fatalf("Error marshalling Set request: %v", err)
	}

	setResponse := httptest.NewRecorder()
	setReq, err := http.NewRequest("POST", "/set", bytes.NewReader(setBody))
	if err != nil {
		t.Fatalf("Error creating Set request: %v", err)
	}

	handler.SetHandler(setResponse, setReq)

	getResponse := httptest.NewRecorder()
	getReq, err := http.NewRequest("GET", "/get?key="+key, nil)
	if err != nil {
		t.Fatalf("Error creating Get request: %v", err)
	}

	handler.GetHandler(getResponse, getReq)

	if getResponse.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, getResponse.Code)
	}

	body, err := ioutil.ReadAll(getResponse.Body)
	if err != nil {
		t.Fatalf("Error reading Get response body: %v", err)
	}

	result := string(body)
	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}
}

func TestAPIDel(t *testing.T) {
	// Your API Del testing code here
}
