package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestShowSnippet(t *testing.T) {
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
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.get(t, tt.urlPath)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestPing(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	// contains a couple of mock loggers (which discard anything written to
	// them).
	app := newTestApplication(t)

	// We then use the httptest.NewTLSServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() to shutdown the server when
	// the test finishes.

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	// The network address that the test server is listening on is contained
	// in the ts.URL field. We can use this along with the ts.Client().Get()
	// method to make a GET /ping request against the test server. This
	// returns a http.Response struct containing the response.

	rs, _, tb := ts.get(t, "/ping")

	// We can then check the value of the response status code and body using
	// the same code as before.
	if rs != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs)
	}

	if string(tb) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

// Unauthenticated users are redirected to the login form.
// Authenticated users are shown the form to create a new snippet.t
func TestCreateSnippetForm(t *testing.T) {
	app := newTestApplication(t)

	//var sr int

	ts := newTestServer(t, app.routes())

	Unauthenticated(t, ts)

	Authenticate(t, ts)

}

func Unauthenticated(t *testing.T, ts *testServer) {

	sc, rh, _ := ts.get(t, "user/create")

	rg := rh.Get("Location")

	if sc != http.StatusSeeOther && rg == "Location: /user/login" {
		t.Errorf("i got: %d and i wanted: %d ", sc, http.StatusOK)
	}

}

func Authenticate(t *testing.T, ts *testServer) error {
	_, _, rb := ts.get(t, "/user/login")

	CSRFtoken := extractCSRFToken(t, rb)

	form := url.Values{}

	form.Add("email", "alice@example.com")
	form.Add("password", "")
	form.Add("csrf_token", CSRFtoken)

	ts.postForm(t, "/user/login", form)

	rs, _, body := ts.get(t, "user/create")

	if rs != http.StatusOK {
		t.Errorf("want %d ; got %d", http.StatusOK, rs)
	}

	formTag := "<form action='/snippet/create' method='POST'>"
	if !bytes.Contains(body, []byte(formTag)) {
		t.Errorf("want body %s to contain %q", body, formTag)
	}
	return nil
}
