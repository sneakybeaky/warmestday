package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	res, err := oc.Client.Get(fmt.Sprintf("%s/data/2.5/onecall?exclude=current,minutely,hourly,alerts&&units=metric&lat=%f&lon=%f&appid=%s", oc.BaseURL, latitude, longitude, oc.AppID))

	if err != nil {
		return weather.Forecast{}, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return weather.Forecast{}, err
	}

	forecast := struct {
		Timezone string `json:"timezone"`
		Daily    []struct {
			Dt   int `json:"dt"`
			Temp struct {
				Max float64 `json:"max"`
			} `json:"temp"`
			Humidity int `json:"humidity"`
		} `json:"daily"`
	}{}

	jsonErr := json.Unmarshal(body, &forecast)
	if jsonErr != nil {
		return weather.Forecast{}, err
	}

	days := make([]weather.Day, len(forecast.Daily))
	for i, day := range forecast.Daily {

		days[i] = weather.Day{
			Date:            time.Unix(int64(day.Dt), 0),
			MaximumTemp:     day.Temp.Max,
			HumidityPercent: day.Humidity,
		}
	}

	return weather.Forecast{
		Timezone: forecast.Timezone,
		Days:     days,
	}, nil
}
