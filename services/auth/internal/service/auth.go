package service

import (
	"context"
	"errors"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	// dummy logic
	if username != "admin" || password != "admin" {
		return "", errors.New("invalid credentials")
	}

	// nanti diganti JWT
	return "dummy-token", nil
}
