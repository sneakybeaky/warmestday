package openweather_test

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"warmestday/weather/openweather"
)

func TestCallIsCorrectlyFormed(t *testing.T) {
	t.Parallel()

	called := false
	wantAppID := "123456789"
	wantLatitude := 1.1
	wantLongitude := -88.76543
	wantUnits := "metric"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		called = true
		// Test request well formed
		if req.URL.Path != "/data/2.5/onecall" {
			t.Errorf("Wanted a path of \"/data/2.5/onecall\" but got $q")
		}

		values := req.URL.Query()
		if gotLatitude := values.Get("lat"); gotLatitude != fmt.Sprintf("%f", wantLatitude) {
			t.Errorf("Wanted a latitude of %f but got %q", wantLatitude, gotLatitude)
		}

		if gotLongitude := values.Get("lon"); gotLongitude != fmt.Sprintf("%f", wantLongitude) {
			t.Errorf("Wanted a longitude of %f but got %q", wantLongitude, gotLongitude)
		}

		if gotAppID := values.Get("appid"); gotAppID != wantAppID {
			t.Errorf("Wanted an appid of %q but got %q", wantAppID, gotAppID)
		}

		if gotUnits := values.Get("units"); gotUnits != wantUnits {
			t.Errorf("Wanted a units of %q but got %q", wantUnits, gotUnits)
		}

		// Send an empty response
		rw.Write([]byte(`{}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	oc := openweather.NewOneCall(wantAppID, func(call *openweather.OneCall) {
		call.Client = server.Client()
		call.BaseURL = server.URL
	})

	_, err := oc.Forecast(wantLatitude, wantLongitude)

	if err != nil {
		t.Fatal(err)
	}

	if called == false {
		t.Fatal("api wasn't called")
	}

}

//go:embed testdata/brighton_forecast.json
var brightonForecast string

func TestResponseConvertedCorrectly(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// Send an empty response
		rw.Write([]byte(brightonForecast))
	}))

	// Close the server when test finishes
	defer server.Close()

	oc := openweather.NewOneCall("1234", func(call *openweather.OneCall) {
		call.Client = server.Client()
		call.BaseURL = server.URL
	})
	got, err := oc.Forecast(1.0, 1.0)

	if err != nil {
		t.Fatal(err)
	}

	wantTimezone := "Europe/London"
	if got.Timezone != wantTimezone {
		t.Errorf("Expected a timezone of %q but got %q", wantTimezone, got.Timezone)
	}

}
