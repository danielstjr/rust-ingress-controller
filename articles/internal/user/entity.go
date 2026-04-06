package user

import "time"

type User struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	DeletedAt *time.Time
}
