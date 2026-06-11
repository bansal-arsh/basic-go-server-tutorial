package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHello(t *testing.T) {
	w := httptest.NewRecorder()
	handleHello(w, nil)

	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Bad response code! Expected %v but received %v.\nBody: %s\n", expectedStatusCode, w.Code, w.Body)
	}

	expectedBody := []byte("Hello, World!\n")
	if !bytes.Equal(w.Body.Bytes(), expectedBody) {
		t.Errorf("Bad response! Expected %q but received %q", expectedBody, w.Body.String())
	}
}

func TestHandleGoodbye(t *testing.T) {
	w := httptest.NewRecorder()
	handleGoodbye(w, nil)

	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Bad response code! Expected %v but received %v.\nBody: %s\n", expectedStatusCode, w.Code, w.Body)
	}

	expectedBody := []byte("Goodbye!\n")
	if !bytes.Equal(w.Body.Bytes(), expectedBody) {
		t.Errorf("Bad response! Expected %q but received %q", expectedBody, w.Body.String())
	}
}
