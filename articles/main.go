package main

import (
	"articles/internal/domain"
	"articles/internal/user"
	"database/sql"
	"encoding/json"
	"errors"
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
	db, err := getDb()
	if err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userController := user.NewController(userService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/user", toHandler(userController.Create))
	r.Get("/user/{userId}", toHandler(userController.Read))
	r.Patch("/user/{userId}", toHandler(userController.Update))
	r.Delete("/user/{userId}", toHandler(userController.Delete))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(notFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(notAllowed)
	})

	http.ListenAndServe(":8080", r)
}

var (
	badRequest    = []byte(`{"message":"bad request"}`)
	internalError = []byte(`{"message":"internal server error"}`)
	notAllowed    = []byte(`{"message":"not allowed"}`)
	notFound      = []byte(`{"message":"not found"}`)
)

type validationErrorResponse struct {
	Message string                 `json:"message"`
	Errors  domain.ValidationError `json:"errors"`
}

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func toHandler(h HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			w.Header().Set("Content-Type", "application/json")

			if errors.Is(err, domain.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write(notFound)
			} else if validationErrors, ok := errors.AsType[domain.ValidationError](err); ok {
				w.WriteHeader(http.StatusNotFound)

				log.Print(validationErrors)

				json.NewEncoder(w).Encode(validationErrorResponse{
					Message: http.StatusText(http.StatusUnprocessableEntity),
					Errors:  validationErrors,
				})
			} else if errors.Is(err, domain.ErrBadRequest) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(badRequest)
			} else {
				log.Print(err)

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(internalError)
			}
		}
	}
}

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
		TZ:       "UTC",
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
