package service

import (
	"context"

	"github.com/afiffaizun/event-driven-backend/services/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
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

	// Ganti comparison langsung dengan bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return "dummy-token", nil // Ganti dengan JWT atau token proper
}

// Tambahkan method untuk register dengan hash
func (s *AuthService) Register(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepo.CreateUser(ctx, username, string(hashedPassword))
}
