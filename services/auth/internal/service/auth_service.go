package service

import (
	"context"
	"errors"
	"time"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/repository"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/security"
)

var ErrInvalidCredential = errors.New("invalid username or password")

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userRepo        repository.UserRepository
	refreshRepo     repository.RefreshTokenRepository
	jwtSecret       string
	refreshDuration time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshRepo repository.RefreshTokenRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		refreshRepo:     refreshRepo,
		jwtSecret:       jwtSecret,
		refreshDuration: 7 * 24 * time.Hour,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	if err := security.CheckPassword(user.Password, password); err != nil {
		return nil, ErrInvalidCredential
	}

	access, err := security.GenerateAccessToken(user.ID, user.Username, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refresh, err := security.GenerateRefreshToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	if err := s.refreshRepo.Save(ctx, &repository.RefreshToken{
		UserID:    user.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(s.refreshDuration),
	}); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	rt, err := s.refreshRepo.FindValid(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	claims, err := security.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil || claims["type"] != "refresh" {
		return "", errors.New("invalid token")
	}

	username, _ := claims["username"].(string)
	return security.GenerateAccessToken(rt.UserID, username, s.jwtSecret)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshRepo.Revoke(ctx, refreshToken)
}
