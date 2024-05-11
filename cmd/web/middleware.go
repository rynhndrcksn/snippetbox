package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

// commonHeaders sets all the default headers we want on every request.
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self';"+
				"style-src 'self' fonts.bunny.net;"+
				"font-src fonts.bunny.net")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// logRequests will log information for each request the site gets.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.RequestURI
		)

		app.logger.Info("Received request", slog.String("ip", ip), slog.String("proto", proto), slog.String("method", method), slog.String("uri", uri))

		next.ServeHTTP(w, r)
	})
}

// recoverPanic checks if there's been a panic and returns a 500 error rather than no response.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic)
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not.
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requireAuthentication makes sure that a user is logged in and authenticated.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
