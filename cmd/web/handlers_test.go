package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"warmestday/cmd/web"
	"warmestday/weather"
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

func TestSummaryHappyPath(t *testing.T) {

	t.Parallel()

	want := "1998-01-11"
	var sf sumarizerFunc = func(_, _ float64) (weather.Summary, error) {
		return weather.Summary{WarmestDay: want}, nil
	}
	app := newTestApplication(withSummarizer(sf))
	ts := newTestServer(app.Routes())

	defer ts.Close()

	got := ts.getWarmestDay(t, 1, 1)

	if got != want {
		t.Fatalf("Wanted %q but got %q", want, got)
	}

}

func TestLocationOutsideEuropeReturnsBadRequest(t *testing.T) {

	t.Parallel()

	var sf sumarizerFunc = func(_, _ float64) (weather.Summary, error) {
		return weather.Summary{}, weather.ErrOutsideEurope
	}
	app := newTestApplication(withSummarizer(sf))
	ts := newTestServer(app.Routes())

	defer ts.Close()

	got, _, _ := ts.get(t, fmt.Sprintf("/summary?lat=1&lon=1"))

	if got != http.StatusBadRequest {
		t.Fatalf("A request for a location outside Europe should return a bad request status, but got %d", got)
	}

}

func TestLatitudeLongitudeDecimalPlaces(t *testing.T) {

	t.Parallel()

	cases := map[string]struct {
		latitude, longitude float64
	}{
		"latitude too many places": {
			latitude: 1.1234567,
		},
		"longitude too many places": {
			longitude: 0.9876543,
		},
	}

	var sf sumarizerFunc = func(_, _ float64) (weather.Summary, error) {
		return weather.Summary{WarmestDay: "want"}, nil
	}
	app := newTestApplication(withSummarizer(sf))
	ts := newTestServer(app.Routes())

	defer ts.Close()

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {

			got, _, _ := ts.get(t, fmt.Sprintf("/summary?lat=%.8f&lon=%.8f", tc.latitude, tc.longitude))

			if got != http.StatusBadRequest {
				t.Fatalf("More than 6 decimal places should return a bad request status, but got %d", got)
			}

		})
	}

}

type sumarizerFunc func(latitude, longitude float64) (weather.Summary, error)

func (f sumarizerFunc) Summarize(latitude, longitude float64) (weather.Summary, error) {
	return f(latitude, longitude)
}

func withSummarizer(summarizer weather.Summarizer) func(application *main.Application) {
	return func(a *main.Application) {
		a.Summarizer = summarizer
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

func (ts *testServer) getWarmestDay(t *testing.T, lat, lon float64) string {
	t.Helper()

	code, _, body := ts.get(t, fmt.Sprintf("/summary?lat=%f&lon=%f", lat, lon))

	if code != http.StatusOK {
		t.Errorf("Wanted a status code of %d but got %d", http.StatusOK, code)
	}

	var summary = struct {
		WarmestDay string
	}{}

	err := json.Unmarshal([]byte(body), &summary)

	if err != nil {
		t.Fatal(err)
	}

	return summary.WarmestDay

}
