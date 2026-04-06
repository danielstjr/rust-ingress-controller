package user

import (
	"articles/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) create(ctx context.Context, name string) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (name) VALUES ($1) RETURNING id, name, created_at, deleted_at",
		name,
	).Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.DeletedAt,
	)

	if err == nil {
		return user, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	return nil, err
}

func (r *repository) findById(ctx context.Context, id uint64) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, created_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.DeletedAt,
	)

	if err == nil {
		return user, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	return nil, err
}

func (r *repository) update(ctx context.Context, user *User) error {
	err := r.db.QueryRowContext(
		ctx,
		"UPDATE users SET name = $1 WHERE id = $2 AND deleted_at IS NULL RETURNING name, created_at, deleted_at",
		user.Name,
		user.ID,
	).Scan(
		&user.Name,
		&user.CreatedAt,
		&user.DeletedAt,
	)

	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}

	return err
}

func (r *repository) deleteById(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return domain.ErrNotFound
	}

	return nil
}
