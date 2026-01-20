package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT id, username, password FROM users WHERE username=$1`,
		username,
	)

	u := &User{}
	if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, username, password string) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO users (username, password) VALUES ($1, $2)`,
		username,
		password,
	)
	return err
}
