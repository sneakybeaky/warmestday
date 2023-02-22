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

/* TODO tests
- no days at all in forecast
*/

func TestWarmestDayIsFirstWhenOnlyOneDayInForecast(t *testing.T) {

	t.Parallel()

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

func TestWarmestDayChosenWhenMoreThanOneDayInForecast(t *testing.T) {

	t.Parallel()

	wantDay := "2020-01-20"

	var f forecastfunc = func(_, _ float64) (forecast.Forecast, error) {
		return forecast.Forecast{
			Timezone: "Europe/London",

			Days: []forecast.Day{{
				Date:        mustTime(t, "2020-01-19"),
				MaximumTemp: 9,
			},
				{
					Date:        mustTime(t, wantDay),
					MaximumTemp: 10,
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

func TestWarmestDayChosenInFirstSevenOnly(t *testing.T) {

	t.Parallel()

	wantDay := "2020-01-20"

	days := make([]forecast.Day, 8)
	days[0] = forecast.Day{
		Date:        mustTime(t, "2020-01-19"),
		MaximumTemp: 8,
	}
	days[6] = forecast.Day{
		Date:        mustTime(t, wantDay),
		MaximumTemp: 9,
	}
	days[7] = forecast.Day{
		Date:        mustTime(t, "2020-01-26"),
		MaximumTemp: 10,
	}

	var f forecastfunc = func(_, _ float64) (forecast.Forecast, error) {
		return forecast.Forecast{
			Timezone: "Europe/London",
			Days:     days,
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

func TestFirstOfManyChosenWhenTheyHaveSameTempAndHumidity(t *testing.T) {

	t.Parallel()

	wantDay := "2020-01-20"

	var f forecastfunc = func(_, _ float64) (forecast.Forecast, error) {
		return forecast.Forecast{
			Timezone: "Europe/London",

			Days: []forecast.Day{{
				Date:            mustTime(t, wantDay),
				MaximumTemp:     9,
				HumidityPercent: 20,
			}, {
				Date:            mustTime(t, wantDay),
				MaximumTemp:     10,
				HumidityPercent: 20,
			}, {
				Date:            mustTime(t, "2020-01-21"),
				MaximumTemp:     10,
				HumidityPercent: 20,
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

func TestDayWithLowestHumidityChosenWhenMoreThanOneDayWithSameTemperature(t *testing.T) {

	t.Parallel()

	wantDay := "2020-01-20"

	var f forecastfunc = func(_, _ float64) (forecast.Forecast, error) {
		return forecast.Forecast{
			Timezone: "Europe/London",

			Days: []forecast.Day{{
				Date:            mustTime(t, "2020-01-19"),
				MaximumTemp:     10,
				HumidityPercent: 11,
			},
				{
					Date:            mustTime(t, wantDay),
					MaximumTemp:     10,
					HumidityPercent: 9,
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
