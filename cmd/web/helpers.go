package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// serverError writes a log entry at Error level and sends a generic 500 http response to the client
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		trace  = string(debug.Stack())
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends the specific http response to the client
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
