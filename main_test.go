package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"reflect"
	"starter-projects/basic-go-server/internal/users"
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

func TestHandleHelloParametrized_Normal(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello?user=TestMan", nil)
	handleHelloParametrized(w, r)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, TestMan!\n"), w, t)
}

func TestHandleHelloParametrized_NoParam(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)
	handleHelloParametrized(w, r)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, User!\n"), w, t)
}

func TestHandleHelloParametrized_WrongParam(t *testing.T) {
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

func TestHandleHelloHeader_Normal(t *testing.T) {
	testManager := users.NewManager()
	testManager.AddUser("foo", "bar", "f.bar@example.com")
	testManager.AddUser("bar", "baz", "b.baz@example.com")
	testManager.AddUser("foo", "baz", "f.baz@example.com")
	testManager.AddUser("baz", "foo", "b.foo@example.com")

	testServer := serverType{manager: testManager}

	req := httptest.NewRequest(http.MethodGet, "/user/hello", nil)
	req.Header.Add("userFirst", "foo")
	req.Header.Add("userLast", "baz")

	w := httptest.NewRecorder()

	testServer.handleHelloHeader(w, req)

	validateCode(http.StatusOK, w, t)
	validateBody([]byte("Hello, foo baz!\nYour email is: f.baz@example.com\n"), w, t)
}

func TestHandleHelloHeader_NoFirstHeader(t *testing.T) {
	testManager := users.NewManager()
	testManager.AddUser("foo", "bar", "f.bar@example.com")
	testManager.AddUser("bar", "baz", "b.baz@example.com")
	testManager.AddUser("foo", "baz", "f.baz@example.com")
	testManager.AddUser("baz", "foo", "b.foo@example.com")

	testServer := serverType{manager: testManager}

	req := httptest.NewRequest(http.MethodGet, "/user/hello", nil)
	req.Header.Add("userLast", "baz")

	w := httptest.NewRecorder()
	testServer.handleHelloHeader(w, req)

	validateCode(http.StatusBadRequest, w, t)
	validateBody([]byte("Invalid first name\n"), w, t)
}

func TestHandleHelloHeader_NoLastHeader(t *testing.T) {
	testManager := users.NewManager()
	testManager.AddUser("foo", "bar", "f.bar@example.com")
	testManager.AddUser("bar", "baz", "b.baz@example.com")
	testManager.AddUser("foo", "baz", "f.baz@example.com")
	testManager.AddUser("baz", "foo", "b.foo@example.com")

	testServer := serverType{manager: testManager}

	req := httptest.NewRequest(http.MethodGet, "/user/hello", nil)
	req.Header.Add("userFirst", "foo")

	w := httptest.NewRecorder()
	testServer.handleHelloHeader(w, req)

	validateCode(http.StatusBadRequest, w, t)
	validateBody([]byte("Invalid last name\n"), w, t)
}

func TestHandleHelloJSON_Normal(t *testing.T) {
	requestStruct := UserData{FirstName: "Test Man"}
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

func TestHandleHelloJSON_EmptyBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/json", nil)
	w := httptest.NewRecorder()
	handleHelloJSON(w, r)

	validateCode(http.StatusBadRequest, w, t)
	validateBody([]byte("Error reading JSON\n"), w, t)
}

func TestHandleHelloJSON_EmptyName(t *testing.T) {
	requestStruct := UserData{FirstName: ""}
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

func TestAddUser_Normal(t *testing.T) {
	testData := UserData{
		FirstName: "Test",
		LastName:  "User Man",
		Email:     "testman@user.com",
	}
	testJSON, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Error marshalling test json: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewBuffer(testJSON))
	r.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testServer := serverType{manager: users.NewManager()}

	testServer.addUser(w, r)
	validateCode(http.StatusCreated, w, t)

	resultUser, err := testServer.manager.GetUserByName(testData.FirstName, testData.LastName)
	if err != nil {
		t.Fatalf("Error retreiving new test user: %v", err)
	}

	resultUserData := convertUserToUserData(resultUser)
	if !reflect.DeepEqual(*resultUserData, testData) {
		t.Fatalf("Retrieved user is not same as created user.\nExpected: %+v\nActual: %+v", testData, resultUserData)
	}
}

func TestAddUser_IncorrectHeader(t *testing.T) {
	testData := UserData{
		FirstName: "Test",
		LastName:  "User Man",
		Email:     "testman@user.com",
	}
	testJSON, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Error marshalling test json: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewBuffer(testJSON))

	w := httptest.NewRecorder()
	testServer := serverType{manager: users.NewManager()}

	testServer.addUser(w, r)
	validateCode(http.StatusUnsupportedMediaType, w, t)
	validateBody([]byte("Unsupported Content-Type header: \"\"\n"), w, t)

	retreivedUser, err := testServer.manager.GetUserByName(testData.FirstName, testData.LastName)
	if err == nil {
		t.Fatalf("User created even though content-type is unsupported: %+v", convertUserToUserData(retreivedUser))
	} else if !errors.Is(err, users.ErrNoResultsFound) {
		t.Fatalf("Error while retreiving user: %v", err)
	}
}

func TestGetUser_Normal(t *testing.T) {
	testManager := users.NewManager()
	testManager.AddUser("foo", "bar", "f.bar@example.com")
	testManager.AddUser("bar", "baz", "b.baz@example.com")
	testManager.AddUser("foo", "baz", "f.baz@example.com")
	testManager.AddUser("baz", "foo", "b.foo@example.com")

	testServer := serverType{manager: testManager}

	testFirstName, testLastName, testEmail := "foo", "baz", "f.baz@example.com"
	testData := struct{ FirstName, LastName string }{testFirstName, testLastName}
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Error marshalling test data: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/get-user", bytes.NewBuffer(jsonBytes))
	r.Header.Add("Content-Type", "application/json")

	testServer.getUser(w, r)
	validateCode(http.StatusOK, w, t)

	responseDecoder := json.NewDecoder(w.Body)
	responseDecoder.DisallowUnknownFields()
	var responseUserData UserData
	err = responseDecoder.Decode(&responseUserData)
	if err != nil {
		t.Fatalf("Error decoding response data: %v", err)
	}

	expectedUserData := UserData{
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmail,
	}
	if !reflect.DeepEqual(expectedUserData, responseUserData) {
		t.Fatalf("Bad response.\nExpected: %+v\nGot: %+v", expectedUserData, responseUserData)
	}
}

func TestGetUser_IncorrectHeader(t *testing.T) {
	testManager := users.NewManager()
	testManager.AddUser("foo", "bar", "f.bar@example.com")
	testManager.AddUser("bar", "baz", "b.baz@example.com")
	testManager.AddUser("foo", "baz", "f.baz@example.com")
	testManager.AddUser("baz", "foo", "b.foo@example.com")

	testServer := serverType{manager: testManager}

	testData := struct{ FirstName, LastName string }{"foo", "baz"}
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Error marshalling test data: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/get-user", bytes.NewBuffer(jsonBytes))
	r.Header.Add("Content-Type", "test")

	testServer.getUser(w, r)
	validateCode(http.StatusUnsupportedMediaType, w, t)
	validateBody([]byte("Unsupported Content-Type header: \"test\"\n"), w, t)
}

func TestGetUser_NoUser(t *testing.T) {
	testManager := users.NewManager()
	testManager.AddUser("foo", "bar", "f.bar@example.com")
	testManager.AddUser("bar", "baz", "b.baz@example.com")
	testManager.AddUser("foo", "baz", "f.baz@example.com")
	testManager.AddUser("baz", "foo", "b.foo@example.com")

	testServer := serverType{manager: testManager}

	testData := struct{ FirstName, LastName string }{"qux", "quz"}
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Error marshalling test data: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/get-user", bytes.NewBuffer(jsonBytes))
	r.Header.Add("Content-Type", "application/json")

	testServer.getUser(w, r)
	validateCode(http.StatusNotFound, w, t)
	validateBody([]byte("No users found\n"), w, t)
}

func TestConvertUserToUserData(t *testing.T) {
	testFirstName := "Test"
	testLastName := "User Man"
	testEmail, err := mail.ParseAddress("testman@user.com")
	if err != nil {
		t.Fatalf("Error while parsing email: %v", err)
	}

	testUser := users.User{
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     *testEmail,
	}
	convertedData := *convertUserToUserData(&testUser)

	expectedData := UserData{
		FirstName: "Test",
		LastName:  "User Man",
		Email:     "testman@user.com",
	}

	if !reflect.DeepEqual(expectedData, convertedData) {
		t.Errorf("Incorrect conversion of User to UserData.\nExpected: %+v\nActual: %+v", expectedData, convertedData)
	}
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
