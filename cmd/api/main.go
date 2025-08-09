package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// global const that store API version
const version = "1.0.0"

// struct that store all config setting for app
type config struct {
	port int
	env  string
}

// struct that hold dependencies for our app
type application struct {
	config config
	logger *log.Logger
}

func main() {
	// create an instance of config
	var cfg config

	// read values from command line flag
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// create logger that logs to the terminal(os.stout)
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// create and instance of the application struct
	app := &application{
		config: cfg,
		logger: logger,
	}

	// create a new router (ServeMux)
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// start the server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)

}
