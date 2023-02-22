package openweather_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"warmestday/weather/openweather"
)

func TestCallIsCorrectlyFormed(t *testing.T) {
	t.Parallel()

	called := false
	wantAppID := "123456789"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		called = true
		// Test request well formed
		if req.URL.Path != "/data/2.5/onecall" {
			t.Errorf("Wanted a path of \"/data/2.5/onecall\" but got $q")
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

	_, err := oc.Forecast(1.0, 2.0)

	if err != nil {
		t.Fatal(err)
	}

	if called == false {
		t.Fatal("api wasn't called")
	}

}
