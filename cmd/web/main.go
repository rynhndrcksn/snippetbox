package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// create a file server which serves files out of the "./ui/static" directory
	// path given is relative to the project directory root
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// use mux.Handle() to register the file server as a handler for all URL paths that start with "/static/"
	// for match paths, we strip the "/static" prefix before the request reaches the file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// register application routes
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
