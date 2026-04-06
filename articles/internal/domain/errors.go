package domain

import (
	"errors"
	"strings"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)

type FieldMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError []FieldMessage

func (e ValidationError) Error() string {
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Message)
	}

	return "Invalid Fields: " + strings.Join(messages, ", ")
}
