package weather

import (
	"errors"
	"fmt"
	"strings"
	"warmestday/forecast"
)

// ErrOutsideEurope signals that the latitude and longitude designate a location outside of Europe
var ErrOutsideEurope = errors.New("location is outside of Europe")

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
// If the location is outside of Europe a ErrOutsideEurope error is returned
func (w Weather) Summary(latitude, longitude float64) (Summary, error) {

	f, err := w.forecaster.Forecast(latitude, longitude)

	if err != nil {
		return Summary{}, err
	}

	continent := strings.Split(f.Timezone, "/")[0]

	if continent != "Europe" {
		return Summary{}, fmt.Errorf("location %q not supported : %w", f.Timezone, ErrOutsideEurope)
	}

	// TODO what if there are no days ?

	days := f.Days

	// We only consider the first 7 days in the forecast
	if len(days) > 7 {
		days = days[:7]
	}

	warmest := days[0]
	for _, day := range days {

		if day.MaximumTemp > warmest.MaximumTemp {
			warmest = day
		}

		if day.MaximumTemp == warmest.MaximumTemp {
			if day.HumidityPercent < warmest.HumidityPercent {
				warmest = day
			}
		}

	}

	return Summary{WarmestDay: warmest.Date.Format("2006-01-02")}, nil
}
