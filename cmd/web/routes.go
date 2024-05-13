package main

import (
	"github.com/justinas/alice"
	"github.com/rynhndrcksn/snippetbox/ui"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create an HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for example, our CSS stylesheet is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Used to determine site is up and running.
	mux.HandleFunc("GET /ping", ping)

	// Create a middleware chain containing our "dynamic" middleware
	// that should only be used on dynamic routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Register application routes - because we're using "dynamic" middleware which returns a http.Handler
	// instead of a http.HandlerFunc, we need to call mux.Handle instead of mux.HandleFunc.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /", dynamic.ThenFunc(app.notFound))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/register", dynamic.ThenFunc(app.userRegister))
	mux.Handle("POST /user/register", dynamic.ThenFunc(app.userRegisterPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated only) routes.
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create a middleware chain containing our "standard" middleware
	// that should be used on every request the webapp receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)

	// If not using alice, middleware must be chained like this.
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
