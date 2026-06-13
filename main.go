package main

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type UserData struct {
	Name string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/goodbye", handleGoodbye)
	mux.HandleFunc("/hello", handleHelloParametrized)
	mux.HandleFunc("/responses/{user}/hello", handleHelloVarUrl)
	mux.HandleFunc("/user/hello", handleHelloHeader)
	mux.HandleFunc("/json", handleHelloJSON)

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

func handleHelloHeader(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("user")
	if username == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

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

	if requestUserData.Name == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	handleHello(w, requestUserData.Name)
}

func handleHello(w http.ResponseWriter, username string) {
	response := "Hello, " + username + "!\n"
	_, err := w.Write([]byte(response))
	if err != nil {
		slog.Error("Error writing response", "err", err)
	}
}
