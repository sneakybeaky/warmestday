package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	"warmestday/forecast/openweather"
	"warmestday/weather"
)

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Weather  weather.Weather
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", ":4000", "HTTP network address")
	appID := flag.String("appid", "", "openweather appid")

	flag.Parse()

	if *appID == "" {
		errorLog.Fatal("You must supply the appid to use")
	}

	weather := weather.NewWeather(openweather.NewOneCall(*appID))

	app := &Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Weather:  weather,
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
