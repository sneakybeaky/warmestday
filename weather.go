package warmestday

import (
	"warmestday/forecast"
)

type Summary struct {
	// // Format YYYY-MM-DD
	WarmestDay string
}

type Weather struct {
	forecaster forecast.Forecaster
}

func NewWeather(forecaster forecast.Forecaster) Weather {
	return Weather{forecaster: forecaster}
}

// Summary returns the warmest day in the next 7 days for the supplied latitude and longitude
func (w Weather) Summary(latitude, longitude float64) (Summary, error) {

	f, err := w.forecaster.Forecast(latitude, longitude)

	if err != nil {
		return Summary{}, err
	}

	return Summary{WarmestDay: f.Days[0].Date.Format("2006-01-02")}, nil
}
