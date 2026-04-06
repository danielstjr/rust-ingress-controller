package user

import (
	"articles/internal/domain"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request) error
	Read(w http.ResponseWriter, r *http.Request) error
	Update(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
}

type controller struct {
	service *Service
}

func NewController(service *Service) Controller {
	return &controller{service: service}
}

type CreateRequest struct {
	Name *string `json:"name"`
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request) error {
	var body CreateRequest

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(&body)
	if err != nil {
		log.Print(err)

		return domain.ErrBadRequest
	}

	if body.Name == nil {
		validations := domain.ValidationError{{Field: "name", Message: "`name` is required"}}

		return validations
	}

	ctx := r.Context()
	user, err := c.service.create(ctx, *body.Name)
	if user != nil && err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		return json.NewEncoder(w).Encode(map[string]any{"id": user.ID, "name": user.Name})
	}

	return err
}

func (c *controller) Read(w http.ResponseWriter, r *http.Request) error {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		log.Print(err)

		return domain.ErrBadRequest
	}

	ctx := r.Context()
	user, err := c.service.get(ctx, userId)
	if user != nil && err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		return json.NewEncoder(w).Encode(map[string]any{"id": user.ID, "name": user.Name})
	}

	return err
}

type UpdateRequest struct {
	Name *string `json:"name"`
}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) error {
	var body UpdateRequest

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(&body)
	if err != nil {
		log.Print(err)

		return domain.ErrBadRequest
	}

	var validations domain.ValidationError

	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		validations = append(validations, domain.FieldMessage{Field: "id", Message: "id must be a uint64"})
	}

	if body.Name == nil {
		validations = append(validations, domain.FieldMessage{Field: "name", Message: "`name` is required"})
	}

	if len(validations) > 0 {
		return validations
	}

	ctx := r.Context()
	user := User{
		ID:   userId,
		Name: *body.Name,
	}

	err = c.service.update(ctx, &user)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		return json.NewEncoder(w).Encode(map[string]any{"id": user.ID, "name": user.Name})
	}

	return err
}

func (c *controller) Delete(w http.ResponseWriter, r *http.Request) error {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		log.Print(err)

		return domain.ErrBadRequest
	}

	ctx := r.Context()
	err = c.service.delete(ctx, userId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	return nil
}
