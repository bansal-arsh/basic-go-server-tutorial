package main

import (
	"log"
	"log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/goodbye", handleGoodbye)
	mux.HandleFunc("/hello", handleHelloParametrized)
	mux.HandleFunc("/responses/{user}/hello", handleHelloVarUrl)

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

	response := "Hello, " + username + "!\n"
	_, err := w.Write([]byte(response))
	if err != nil {
		slog.Error("Error writing response", "err", err)
		return
	}
}

func handleHelloVarUrl(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")
	response := "Hello, " + username + "!\n"
	_, err := w.Write([]byte(response))
	if err != nil {
		slog.Error("Error writing response", "err", err)
	}
}
