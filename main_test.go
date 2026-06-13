package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	w := httptest.NewRecorder()
	handleRoot(w, nil)

	expectedStatusCode := http.StatusOK
	validateCode(expectedStatusCode, w, t)

	expectedBody := []byte("Welcome to the Homepage!\n")
	validateBody(expectedBody, w, t)
}

func TestHandleGoodbye(t *testing.T) {
	w := httptest.NewRecorder()
	handleGoodbye(w, nil)

	expectedStatusCode := http.StatusOK
	validateCode(expectedStatusCode, w, t)

	expectedBody := []byte("Goodbye!\n")
	validateBody(expectedBody, w, t)
}

func TestHandleHelloParametrizedNormal(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello?user=TestMan", nil)
	handleHelloParametrized(w, r)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, TestMan!\n"), w, t)
}

func TestHandleHelloParametrizedNoParam(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)
	handleHelloParametrized(w, r)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, User!\n"), w, t)
}

func TestHandleHelloParametrizedWrongParam(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello?foo=bar", nil)
	handleHelloParametrized(w, r)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, User!\n"), w, t)
}

func TestHandleHelloVarUrl(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/responses/TestMan/hello", nil)
	req.SetPathValue("user", "TestMan")

	w := httptest.NewRecorder()

	handleHelloVarUrl(w, req)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, TestMan!\n"), w, t)
}

func TestHandleHelloHeaderNormal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello", nil)
	req.Header.Add("user", "Test Man")

	w := httptest.NewRecorder()

	handleHelloHeader(w, req)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, Test Man!\n"), w, t)
}

func TestHandleHelloHeaderNoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello", nil)

	w := httptest.NewRecorder()

	handleHelloHeader(w, req)

	validateCode(http.StatusBadRequest, w, t)
	validateBody([]byte("Invalid username\n"), w, t)
}

func TestHandleHelloJSON(t *testing.T) {
	requestStruct := UserData{Name: "Test Man"}
	requestData, err := json.Marshal(requestStruct)
	if err != nil {
		t.Fatalf("Error marshalling test struct: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(requestData))
	w := httptest.NewRecorder()
	handleHelloJSON(w, r)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, Test Man!\n"), w, t)
}

func TestHandleHelloJSONEmptyBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/json", nil)
	w := httptest.NewRecorder()
	handleHelloJSON(w, r)

	validateCode(http.StatusBadRequest, w, t)
	validateBody([]byte("Error reading JSON\n"), w, t)
}

func TestHandleHelloJSONEmptyName(t *testing.T) {
	requestStruct := UserData{Name: ""}
	requestData, err := json.Marshal(requestStruct)
	if err != nil {
		t.Fatalf("Error marshalling test struct: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(requestData))
	w := httptest.NewRecorder()
	handleHelloJSON(w, r)

	validateCode(http.StatusBadRequest, w, t)
	validateBody([]byte("Invalid username\n"), w, t)
}

func validateCode(expectedStatusCode int, w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != expectedStatusCode {
		t.Errorf("Bad response code! Expected %v but received %v.\nBody: %s\n", expectedStatusCode, w.Code, w.Body)
	}
}

func validateBody(expectedBody []byte, w *httptest.ResponseRecorder, t *testing.T) {
	if !bytes.Equal(w.Body.Bytes(), expectedBody) {
		t.Errorf("Bad response! Expected %q but received %q", expectedBody, w.Body.String())
	}
}
