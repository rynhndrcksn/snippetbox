package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rynhndrcksn/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// httptest.NewRecorder will record the response status code, headers, and body
	// instead of writing them to an HTTP connection.
	rr := httptest.NewRecorder()

	// Initialize a new dummy http.Request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the ping handler function, passing in the
	// httptest.ResponseRecorder and http.Request
	ping(rr, r)

	// Call Result() method to get the http.Response
	rs := rr.Result()

	// Check that the status code was 200
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// Check body response is "OK".
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
