package auth

import "immodi/novel-site/internal/app/services/auth"

type AuthHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(service auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}
