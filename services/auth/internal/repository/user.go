package repository

import "context"

type User struct {
	ID       int64
	Username string
	Password string
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
}
