package openweather

import (
	"net/http"
	"warmestday/weather"
)

// OneCall is a simple client for the openweather api
// See https://openweathermap.org/api/one-call-api for details
type OneCall struct {
	AppID   string
	Client  *http.Client
	BaseURL string
}

func NewOneCall(appID string, opts ...func(call *OneCall)) OneCall {
	oc := OneCall{
		AppID:   appID,
		Client:  http.DefaultClient,
		BaseURL: "https://api.openweathermap.org",
	}

	for _, opt := range opts {
		opt(&oc)
	}

	return oc
}

func (oc OneCall) Forecast(latitude, longitude float32) (weather.Forecast, error) {
	_, _ = oc.Client.Get(oc.BaseURL + "/data/2.5/onecall")
	return weather.Forecast{}, nil
}
