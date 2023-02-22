package openweather

import (
	"fmt"
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

func (oc OneCall) Forecast(latitude, longitude float64) (weather.Forecast, error) {
	_, _ = oc.Client.Get(fmt.Sprintf("%s/data/2.5/onecall?lat=%f&lon=%f", oc.BaseURL, latitude, longitude))
	return weather.Forecast{}, nil
}
