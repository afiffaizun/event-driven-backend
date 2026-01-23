package service

import (
	"context"
	"errors"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/repository"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/security"
)

var ErrInvalidCredential = errors.New("invalid username or password")

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

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", ErrInvalidCredential
	}

	if err := security.CheckPassword(user.Password, password); err != nil {
		return "", ErrInvalidCredential
	}

	return security.GenerateToken(user.ID, user.Username, s.jwtSecret)
}
