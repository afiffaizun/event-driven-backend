package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRefreshTokenRepository(db *pgxpool.Pool) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (r *PostgresRefreshTokenRepository) Save(ctx context.Context, t *RefreshToken) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, t.UserID, t.Token, t.ExpiresAt)
	return err
}

func (r *PostgresRefreshTokenRepository) FindValid(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, token, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token = $1
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
	`, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *PostgresRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token = $1
	`, token)
	return err
}

func (r *PostgresRefreshTokenRepository) RevokeAllByUser(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1
	`, userID)
	return err
}
