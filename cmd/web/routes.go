package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes(staticDir string) http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Path given is relative to the project directory root.
	fileServer := http.FileServer(http.Dir(staticDir))

	// Use mux.Handle() to register the file server as a handler for all URL paths that start with "/static/"
	// For matching paths, we strip the "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register application routes.
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our "standard" middleware
	// that should be used on every request the webapp receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)

	// If not using alice, middleware must be chained like this.
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
