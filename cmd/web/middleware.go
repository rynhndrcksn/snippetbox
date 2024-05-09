package main

import "net/http"

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

//https://fonts.bunny.net
