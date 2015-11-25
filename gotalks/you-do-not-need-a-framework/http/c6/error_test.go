package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func (e *HTTPError) shouldMatchJSON(expects interface{}, t *testing.T) {
	expected, err := json.Marshal(expects)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	if err := e.JSON(recorder); err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(recorder.Body.Bytes(), expected) != 0 {
		t.Fatalf("does not match")
	}
}

func TestShouldMatch(t *testing.T) {
	expecting := struct {
		Status    int    `json:"status"`
		Error     string `json:"error"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}{
		Status:    400,
		Error:     http.StatusText(400),
		Message:   "bad error",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	err := NewHTTPError(400, "bad error")
	err.shouldMatchJSON(expecting, t)
}
