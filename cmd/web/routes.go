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

	// Create a middleware chain containing our "dynamic" middleware
	// that should only be used on dynamic routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Register application routes - because we're using "dynamic" middleware which returns a http.Handler
	// instead of a http.HandlerFunc, we need to call mux.Handle instead of mux.HandleFunc.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /", dynamic.ThenFunc(app.notFound))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /user/register", dynamic.ThenFunc(app.userRegister))
	mux.Handle("POST /user/register", dynamic.ThenFunc(app.userRegisterPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// Create a middleware chain containing our "standard" middleware
	// that should be used on every request the webapp receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)

	// If not using alice, middleware must be chained like this.
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
