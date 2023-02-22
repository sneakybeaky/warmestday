package weather

import "time"

// Day holds forecast data for a day.
type Day struct {

	// Date is a point in the day for this forecast.
	Date time.Time

	// MaximumTemp is the maximum temperature for this day in degrees celsius
	MaximumTemp float64

	// HumidityPercent is the humidity percentage
	HumidityPercent int
}

// Forecast holds a locations timezone and forecast for a number of days in the future.
type Forecast struct {

	// Timezone is in the IANA tz database format, e.g. Europe/London
	Timezone string
	Days     []Day
}

// Forecaster returns the forecast for a supplied location.
type Forecaster interface {
	Forecast(latitude, longitude float64) (Forecast, error)
}
