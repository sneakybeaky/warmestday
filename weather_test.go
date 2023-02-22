package warmestday_test

import (
	"testing"
	"time"
	"warmestday"
	"warmestday/forecast"
)

type forecastfunc func(latitude, longitude float64) (forecast.Forecast, error)

func (f forecastfunc) Forecast(latitude, longitude float64) (forecast.Forecast, error) {
	return f(latitude, longitude)
}

func TestWarmestDayIsFirstWhenOnlyOneDayInForecast(t *testing.T) {

	wantDay := "2020-01-20"

	var f forecastfunc = func(_, _ float64) (forecast.Forecast, error) {
		return forecast.Forecast{
			Timezone: "Europe/London",
			Days: []forecast.Day{{
				Date: mustTime(t, wantDay),
			}},
		}, nil
	}

	w := warmestday.NewWeather(f)
	got, err := w.Summary(0, 0)

	if err != nil {
		t.Fatal(err)
	}

	if got.WarmestDay != wantDay {
		t.Fatalf("Expected %q but got %q", wantDay, got.WarmestDay)
	}
}

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", value)

	if err != nil {
		t.Fatal(err)
	}
	return parsed

}
