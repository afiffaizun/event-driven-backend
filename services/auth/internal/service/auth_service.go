package service

import (
	"context"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/repository"
	"github.com/afiffaizun/event-driven-backend/services/auth/internal/security"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := security.CheckPassword(user.Password, password); err != nil {
		return "", ErrInvalidCredentials
	}

	// nanti diganti JWT
	return "dummy-token", nil
}
