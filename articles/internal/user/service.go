package user

import (
	"articles/internal/domain"
	"context"
)

type Repository interface {
	create(ctx context.Context, name string) (*User, error)
	findById(ctx context.Context, id uint64) (*User, error)
	update(ctx context.Context, user *User) error
	deleteById(ctx context.Context, id uint64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) create(ctx context.Context, name string) (*User, error) {
	var errors domain.ValidationError
	if len(name) == 0 {
		errors = append(errors, domain.FieldMessage{Field: "name", Message: "Cannot be empty"})
	} else if len(name) > 256 {
		errors = append(errors, domain.FieldMessage{Field: "name", Message: "Must be less than 257 characters"})
	}

	if len(errors) > 0 {
		return nil, errors
	}

	user, err := s.repo.create(ctx, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) get(ctx context.Context, id uint64) (*User, error) {
	if id == 0 {
		return nil, domain.ValidationError{domain.FieldMessage{Field: "id", Message: "Invalid id provided"}}
	}

	return s.repo.findById(ctx, id)
}

func (s *Service) update(ctx context.Context, user *User) error {
	if user == nil {
		return domain.ErrNotFound
	}

	var errors domain.ValidationError
	if user.ID == 0 {
		errors = append(errors, domain.FieldMessage{Field: "id", Message: "Invalid id provided"})
	}

	if len(user.Name) == 0 {
		errors = append(errors, domain.FieldMessage{Field: "name", Message: "Cannot be empty"})
	} else if len(user.Name) > 256 {
		errors = append(errors, domain.FieldMessage{Field: "name", Message: "Must be less than 257 characters"})
	}

	if len(errors) > 0 {
		return errors
	}

	return s.repo.update(ctx, user)
}

func (s *Service) delete(ctx context.Context, id uint64) error {
	return s.repo.deleteById(ctx, id)
}
