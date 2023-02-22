package main_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"warmestday/cmd/web"
)

func TestPing(t *testing.T) {

	t.Parallel()
	app := newTestApplication()

	ts := newTestServer(app.Routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("Wanted a status code of %d but got %d", http.StatusOK, code)
	}

	if body != "OK" {
		t.Errorf("Wanted a body of \"OK\" but got %q", body)
	}

}

func newTestApplication(opts ...func(application *main.Application)) *main.Application {

	app := &main.Application{
		ErrorLog: log.New(io.Discard, "", 0),
		InfoLog:  log.New(io.Discard, "", 0),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app

}

type testServer struct {
	*httptest.Server
}

func newTestServer(h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	t.Helper()

	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
