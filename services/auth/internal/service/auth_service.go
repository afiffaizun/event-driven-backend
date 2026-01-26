package service

import (
	"context"
	"errors"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/repository"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/security"
)

var ErrInvalidCredential = errors.New("invalid username or password")

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
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

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
	claims, err := security.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return "", err
	}

	if claims["type"] != "refresh" {
		return "", errors.New("invalid token type")
	}

	userID := int64(claims["sub"].(float64))
	username := claims["username"]

	return security.GenerateAccessToken(userID, username.(string), s.jwtSecret)
}
