package openweather_test

import (
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

		// Send response to be tested
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
