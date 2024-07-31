package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHelloHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloHandler)

	originalHostname := os.Hostname
	defer func() { os.Hostname = originalHostname }()
	os.Hostname = func() (string, error) { return "test-hostname", nil }

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := Response{
		Timestamp: time.Now().Format(time.RFC3339),
		Hostname:  "test-hostname",
	}

	var actual Response
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if actual.Hostname != expected.Hostname {
		t.Errorf("handler returned unexpected hostname: got %v want %v", actual.Hostname, expected.Hostname)
	}

	if _, err := time.Parse(time.RFC3339, actual.Timestamp); err != nil {
		t.Errorf("handler returned invalid timestamp: %v", actual.Timestamp)
	}
}
