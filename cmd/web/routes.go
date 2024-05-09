package main

import "net/http"

func (app *application) routes(staticDir string) http.Handler {
	mux := http.NewServeMux()

	// create a file server which serves files out of the "./ui/static" directory
	// path given is relative to the project directory root
	fileServer := http.FileServer(http.Dir(staticDir))

	// use mux.Handle() to register the file server as a handler for all URL paths that start with "/static/"
	// for match paths, we strip the "/static" prefix before the request reaches the file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// register application routes
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return commonHeaders(mux)
}
