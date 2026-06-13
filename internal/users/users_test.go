package users

import (
	"errors"
	"net/mail"
	"reflect"
	"testing"
)

func TestAddUserNormal(t *testing.T) {
	firstName := "Test"
	lastName := "User Man"
	email, err := mail.ParseAddress("testman@user.com")
	if err != nil {
		t.Fatalf("Error parsing test email address. Error: %v", err)
	}

	testManagerPtr := NewManager()
	err = testManagerPtr.AddUser(firstName, lastName, email.String())
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	usersListLength := len(testManagerPtr.users)
	if usersListLength < 1 {
		t.Fatalf("No user created")
	} else if usersListLength > 1 {
		t.Fatalf("Too many (%v) users created", usersListLength)
	}

	expectedUser := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     *email,
	}
	foundUser := testManagerPtr.users[0]
	if !reflect.DeepEqual(expectedUser, foundUser) {
		t.Fatalf("Bad user data.\nExpected: %+v\nReceived: %+v", expectedUser, foundUser)
	}
}

func TestAddUserInvalidEmail(t *testing.T) {
	firstName := "Test"
	lastName := "User Man"
	email := "testuserman"

	testManager := NewManager()
	err := testManager.AddUser(firstName, lastName, email)
	if err == nil {
		t.Fatalf("Error expected but no error received for invalid email")
	}

	expectedErrMsg := "Invalid email address: testuserman"
	if err.Error() != expectedErrMsg {
		t.Fatalf("Bad error message. Expected %q but received %q", expectedErrMsg, err.Error())
	}

	if len(testManager.users) > 0 {
		t.Fatalf("%d users created even though error returned", len(testManager.users))
	}
}

func TestAddUserEmptyFirstName(t *testing.T) {
	firstName := ""
	lastName := "User Man"
	email, err := mail.ParseAddress("testman@user.com")
	if err != nil {
		t.Fatalf("Error parsing test address: %v", err)
	}

	testManager := NewManager()
	err = testManager.AddUser(firstName, lastName, email.String())
	if err == nil {
		t.Fatalf("Expected error but no error received")
	}

	expectedErrMsg := "Invalid first name: \"\""
	if err.Error() != expectedErrMsg {
		t.Fatalf("Bad error message. Expected %q but received %q", expectedErrMsg, err.Error())
	}

	if len(testManager.users) > 0 {
		t.Fatalf("%d users created even though error returned", len(testManager.users))
	}
}

func TestAddUserEmptyLastName(t *testing.T) {
	firstName := "Test"
	lastName := ""
	email, err := mail.ParseAddress("testman@user.com")
	if err != nil {
		t.Fatalf("Error parsing test address: %v", err)
	}

	testManager := NewManager()
	err = testManager.AddUser(firstName, lastName, email.String())
	if err == nil {
		t.Fatalf("Expected error but no error received")
	}

	expectedErrMsg := "Invalid last name: \"\""
	if err.Error() != expectedErrMsg {
		t.Fatalf("Bad error message. Expected %q but received %q", expectedErrMsg, err.Error())
	}

	if len(testManager.users) > 0 {
		t.Fatalf("%d users created even though error returned", len(testManager.users))
	}
}

func TestUserDuplicateName(t *testing.T) {
	firstName := "Test"
	lastName := "User Man"
	email, err := mail.ParseAddress("testman@user.com")
	if err != nil {
		t.Fatalf("Error creating test email: %v", err)
	}

	testManager := NewManager()
	err = testManager.AddUser(firstName, lastName, email.String())
	if err != nil {
		t.Fatalf("Error adding first user: %v", err)
	}

	err = testManager.AddUser(firstName, lastName, email.String())
	if err == nil {
		t.Fatalf("Error expected but no error received")
	}

	expectedError := "User already exists"
	if err.Error() != expectedError {
		t.Fatalf("Bad error message. Expected %q but received %q", expectedError, err.Error())
	}

	if len(testManager.users) != 1 {
		t.Fatalf("Expected 1 user but %d users created", len(testManager.users))
	}
}

func TestGetUserByName(t *testing.T) {
	testManager := NewManager()

	err := testManager.AddUser("foo", "bar", "f.bar@example.com")
	if err != nil {
		t.Fatalf("Error adding user: %v", err)
	}
	err = testManager.AddUser("bar", "baz", "b.baz@example.com")
	if err != nil {
		t.Fatalf("Error adding user: %v", err)
	}
	err = testManager.AddUser("foo", "baz", "f.baz@example.com")
	if err != nil {
		t.Fatalf("Error adding user: %v", err)
	}
	err = testManager.AddUser("baz", "foo", "b.foo@example.com")
	if err != nil {
		t.Fatalf("Error adding user: %v", err)
	}

	tests := map[string]struct {
		firstName    string
		lastName     string
		expectedUser *User
		expectedErr  error
	}{
		"Normal lookup": {
			firstName:    "bar",
			lastName:     "baz",
			expectedUser: &testManager.users[1],
			expectedErr:  nil,
		},
		"Last element lookup": {
			firstName:    "baz",
			lastName:     "foo",
			expectedUser: &testManager.users[3],
			expectedErr:  nil,
		},
		"No match lookup": {
			firstName:    "qux",
			lastName:     "quz",
			expectedUser: nil,
			expectedErr:  ErrNoResultsFound,
		},
		"Partial match lookup": {
			firstName:    "foo",
			lastName:     "foo",
			expectedUser: nil,
			expectedErr:  ErrNoResultsFound,
		},
		"Empty first name": {
			firstName:    "",
			lastName:     "foo",
			expectedUser: nil,
			expectedErr:  ErrNoResultsFound,
		},
		"Empty last name": {
			firstName:    "foo",
			lastName:     "",
			expectedUser: nil,
			expectedErr:  ErrNoResultsFound,
		},
	}

	for testName, test := range tests {
		responseUser, err := testManager.GetUserByName(test.firstName, test.lastName)
		if !reflect.DeepEqual(responseUser, test.expectedUser) {
			t.Errorf("%s: Invalid user struct.\nExpected: %+v\nReceived: %+v", testName, test.expectedUser, responseUser)
		}
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: Invalid error.\nExpected: %v\nReceived: %v", testName, test.expectedErr, err)
		}
	}
}
