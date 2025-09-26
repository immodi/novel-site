package admin

import (
	"immodi/novel-site/internal/app/services/auth"
	"immodi/novel-site/internal/app/services/profile"
)

type AdminHandler struct {
	authService    auth.AuthService
	profileService profile.ProfileService
}

func NewAuthHandler(service auth.AuthService, profileService profile.ProfileService) *AdminHandler {
	return &AdminHandler{authService: service, profileService: profileService}
}
