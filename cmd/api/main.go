package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// global const that store API version
const version = "1.0.0"

// struct that store all config setting for app
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
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

	// load the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — continuing with environment and flags")
	}

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	// create logger that logs to the terminal(os.stout)
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// Open database connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// Close pool when main() ends
	defer db.Close()

	logger.Printf("database connection pool established")

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
	err = srv.ListenAndServe()
	logger.Fatal(err)

}

// openDB() returns a sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	//Set the max open (in-use + idle) connections
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the max idle connections
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// parse idle timeout string into a time.Duration
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// set the max idle time
	db.SetConnMaxIdleTime(duration)

	// create context with 5 second time out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try connecting; fails if not successful within 5 seconds
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
