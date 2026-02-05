package repository

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, t *RefreshToken) error
	FindValid(ctx context.Context, token string) (*RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllByUser(ctx context.Context, userID int64) error
}
