package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"warmestday/weather"
)

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

func (app *Application) summary(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		app.ErrorLog.Printf("Unable to parse latitude : %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		app.ErrorLog.Printf("Unable to parse longitude : %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	summary, err := app.Summarizer.Summarize(lat, lon)

	if err != nil {

		switch {
		case errors.Is(err, weather.ErrOutsideEurope):
			app.ErrorLog.Printf("Location outside of Europe : %v", err)
			http.Error(w, "Location must be inside Europe", http.StatusBadRequest)
			return

		default:
			app.ErrorLog.Printf("Unable to get summary : %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	body, err := json.Marshal(summary)

	if err != nil {
		app.ErrorLog.Printf("Unable to encode summary : %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)

}
