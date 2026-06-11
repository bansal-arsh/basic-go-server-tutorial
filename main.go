package main

import (
	"bytes"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/goodbye", handleGoodbye)
	mux.HandleFunc("/hello", handleHelloParametrized)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello, World!\n"))
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

	var output bytes.Buffer
	output.WriteString("Hello, ")
	output.WriteString(username)
	output.WriteString("!\n")

	_, err := w.Write(output.Bytes())
	if err != nil {
		slog.Error("Error writing response", "err", err)
		return
	}
}
