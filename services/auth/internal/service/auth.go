package service

import (
	"context"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username != "admin" || password != "admin" {
		return "", ErrInvalidCredentials
	}
	return "dummy-token", nil
}

