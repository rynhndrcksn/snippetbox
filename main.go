package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippetbox"))
}

func main() {
	// make a new server mux and register "/" to the home function
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// print a log message saying the server started
	log.Print("Starting server on :8080")

	// start our server and log any fatal errors that get returned
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
