package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Displays a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}

func main() {
	// make a new server mux and register "/" to the home function
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home) // {$} is a special character sequence that restricts "/" to only "/" instead of being a wildcard
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// print a log message saying the server started
	log.Print("Starting server on :8080")

	// start our server and log any fatal errors that get returned
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
