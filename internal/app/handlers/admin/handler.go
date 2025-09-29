package admin

import (
	"immodi/novel-site/internal/app/services/auth"
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/app/services/profile"
)

type AdminHandler struct {
	authService    auth.AuthService
	indexService   index.HomeService
	profileService profile.ProfileService
}

func NewAuthHandler(service auth.AuthService, profileService profile.ProfileService, indexService index.HomeService) *AdminHandler {
	return &AdminHandler{authService: service, profileService: profileService, indexService: indexService}
}
