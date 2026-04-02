package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lib/pq"
)

func main() {
	status := "Success"
	db, err := getDb()
	if err == nil {
		defer db.Close()
	} else {
		log.Print(err)
		status = "Failure"
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Basically a health endpoint for base project infra, will be moved to actual /health later
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		output := struct {
			DbLoaded string `json:"db_loaded"`
		}{status}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output)
	})

	http.ListenAndServe(":8080", r)
}

/**
 * Simple docker pgsql hookup
 */
func getDb() (*sql.DB, error) {
	port, err := strconv.ParseUint(os.Getenv("DB_PORT"), 10, 16)
	if err != nil {
		return nil, err
	}

	config := pq.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     uint16(port),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  pq.SSLModeDisable,
	}

	connection, err := pq.NewConnectorConfig(config)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connection)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	return db, nil
}
