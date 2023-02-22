package openweather_test

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"warmestday/weather"
	"warmestday/weather/openweather"

	"github.com/google/go-cmp/cmp"
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

	// TODO load from golden file maybe
	want := weather.Forecast{
		Timezone: "Europe/London",
		Days: []weather.Day{
			{
				Date:            mustTime(t, "2023-02-22T12:00:00Z"),
				MaximumTemp:     8.25,
				HumidityPercent: 94,
			},
			{
				Date:            mustTime(t, "2023-02-23T12:00:00Z"),
				MaximumTemp:     7.24,
				HumidityPercent: 85,
			},
			{
				Date:            mustTime(t, "2023-02-24T12:00:00Z"),
				MaximumTemp:     7.13,
				HumidityPercent: 87,
			},
			{
				Date:            mustTime(t, "2023-02-25T12:00:00Z"),
				MaximumTemp:     6.48,
				HumidityPercent: 59,
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("response mismatch (-want +got):\n%s", diff)
	}

}

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02T15:04:05Z07", value)

	if err != nil {
		t.Fatal(err)
	}

	return parsed
}
