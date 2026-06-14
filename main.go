package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"starter-projects/basic-go-server/internal/users"
)

type UserData struct {
	FirstName string
	LastName  string
	Email     string
}

type serverType struct {
	manager *users.Manager
}

func main() {
	server := serverType{manager: users.NewManager()}
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/goodbye/", handleGoodbye)
	mux.HandleFunc("/hello/", handleHelloParametrized)
	mux.HandleFunc("/responses/{user}/hello/", handleHelloVarUrl)
	mux.HandleFunc("/user/hello/", server.handleHelloHeader)
	mux.HandleFunc("POST /json", handleHelloJSON)
	mux.HandleFunc("POST /add-user", server.addUser)
	mux.HandleFunc("POST /get-user", server.getUser)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Welcome to the Homepage!\n"))
	if err != nil {
		slog.Error("Error writing response", "err", err)
		return
	}
}

func handleGoodbye(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Goodbye!\n"))
	if err != nil {
		slog.Error("Error writing response", "err", err)
		return
	}
}

func handleHelloParametrized(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	userList, ok := params["user"]
	var username string
	if ok {
		username = userList[0]
	} else {
		username = "User"
	}

	handleHello(w, username)
}

func handleHelloVarUrl(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")
	handleHello(w, username)
}

func handleHelloJSON(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil || len(jsonBytes) < 1 {
		slog.Error("Error reading JSON bytes from body", "err", err)
		http.Error(w, "Error reading JSON", http.StatusBadRequest)
		return
	}

	var requestUserData UserData
	err = json.Unmarshal(jsonBytes, &requestUserData)
	if err != nil {
		slog.Error("Error unmarshalling JSON", "err", err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	if requestUserData.FirstName == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	handleHello(w, requestUserData.FirstName)
}

func handleHello(w http.ResponseWriter, username string) {
	response := "Hello, " + username + "!\n"
	_, err := w.Write([]byte(response))
	if err != nil {
		slog.Error("Error writing response", "err", err)
	}
}

func (s *serverType) handleHelloHeader(w http.ResponseWriter, r *http.Request) {
	firstName := r.Header.Get("userFirst")
	if firstName == "" {
		http.Error(w, "Invalid first name", http.StatusBadRequest)
		return
	}

	lastName := r.Header.Get("userLast")
	if lastName == "" {
		http.Error(w, "Invalid last name", http.StatusBadRequest)
		return
	}

	user, err := s.manager.GetUserByName(firstName, lastName)
	if err != nil {
		if errors.Is(err, users.ErrNoResultsFound) {
			http.Error(w, "No users found", http.StatusNotFound)
		} else {
			slog.Error("Error retreiving user", "err", err)
			http.Error(w, "Error retreiving user", http.StatusBadRequest)
		}
		return
	}

	userData := convertUserToUserData(user)
	response := fmt.Sprintf("Hello, %s %s!\nYour email is: %s\n", firstName, lastName, userData.Email)
	w.Write([]byte(response))
}

func (s *serverType) addUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errMsg := fmt.Sprintf("Unsupported Content-Type header: %q", contentType)
		http.Error(w, errMsg, http.StatusUnsupportedMediaType)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	var newUserData UserData
	err := jsonDecoder.Decode(&newUserData)
	if err != nil {
		slog.Error("Error unmarshalling user data", "err", err)
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	err = s.manager.AddUser(newUserData.FirstName, newUserData.LastName, newUserData.Email)
	if err != nil {
		slog.Error("Error while adding user", "err", err)
		http.Error(w, "Error adding user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *serverType) getUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errMsg := fmt.Sprintf("Unsupported Content-Type header: %q", contentType)
		http.Error(w, errMsg, http.StatusUnsupportedMediaType)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	var requestUserName struct{ FirstName, LastName string }
	err := jsonDecoder.Decode(&requestUserName)
	if err != nil {
		slog.Error("Error unmarshalling request body", "err", err)
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	user, err := s.manager.GetUserByName(requestUserName.FirstName, requestUserName.LastName)
	if err != nil {
		if errors.Is(err, users.ErrNoResultsFound) {
			http.Error(w, "No users found", http.StatusNotFound)
		} else {
			slog.Error("Error retreiving user", "err", err)
			http.Error(w, "Error retreiving user", http.StatusBadRequest)
		}
		return
	}

	userData := convertUserToUserData(user)
	responseBytes, err := json.Marshal(userData)
	if err != nil {
		slog.Error("Error marshalling response user data", "err", err)
		http.Error(w, "Error while sending data", http.StatusInternalServerError)
		return
	}

	w.Write(responseBytes)
}

func convertUserToUserData(u *users.User) *UserData {
	return &UserData{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email.Address,
	}
}
