package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/rynhndrcksn/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "Snippet content",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running an end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/register request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/register")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action=\"/user/register\" method=\"post\" novalidate>"
	)

	tests := []struct {
		name                string
		userName            string
		userEmail           string
		userPassword        string
		userConfirmPassword string
		csrfToken           string
		wantCode            int
		wantFormTag         string
	}{
		{
			name:                "Valid submission",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: validPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusSeeOther,
		},
		{
			name:                "Invalid CSRF Token",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: validPassword,
			csrfToken:           "wrongToken",
			wantCode:            http.StatusBadRequest,
		},
		{
			name:                "Empty name",
			userName:            "",
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: validPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Empty email",
			userName:            validName,
			userEmail:           "",
			userPassword:        validPassword,
			userConfirmPassword: validPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Empty password",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        "",
			userConfirmPassword: "",
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Invalid email",
			userName:            validName,
			userEmail:           "bob@example.",
			userPassword:        validPassword,
			userConfirmPassword: validPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Short password",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        "pa$$",
			userConfirmPassword: "pa$$",
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Passwords don't match",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        "Password1!",
			userConfirmPassword: "Password1",
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Duplicate email",
			userName:            validName,
			userEmail:           "dupe@example.com",
			userPassword:        validPassword,
			userConfirmPassword: validPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("confirm-password", tt.userConfirmPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/register", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}
