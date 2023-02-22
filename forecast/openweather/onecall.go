package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"warmestday/forecast"
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

func (oc OneCall) Forecast(latitude, longitude float64) (forecast.Forecast, error) {
	res, err := oc.Client.Get(fmt.Sprintf("%s/data/2.5/onecall?exclude=current,minutely,hourly,alerts&&units=metric&lat=%f&lon=%f&appid=%s", oc.BaseURL, latitude, longitude, oc.AppID))

	if err != nil {
		return forecast.Forecast{}, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return forecast.Forecast{}, err
	}

	got := struct {
		Timezone string `json:"timezone"`
		Daily    []struct {
			Dt   int `json:"dt"`
			Temp struct {
				Max float64 `json:"max"`
			} `json:"temp"`
			Humidity int `json:"humidity"`
		} `json:"daily"`
	}{}

	err = json.Unmarshal(body, &got)
	if err != nil {
		return forecast.Forecast{}, err
	}

	days := make([]forecast.Day, len(got.Daily))
	for i, day := range got.Daily {

		days[i] = forecast.Day{
			Date:            time.Unix(int64(day.Dt), 0),
			MaximumTemp:     day.Temp.Max,
			HumidityPercent: day.Humidity,
		}
	}

	return forecast.Forecast{
		Timezone: got.Timezone,
		Days:     days,
	}, nil
}
