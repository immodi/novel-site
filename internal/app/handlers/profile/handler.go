package profile

import (
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/app/services/profile"
)

type ProfileHandler struct {
	profileService profile.ProfileService
	homeService    index.HomeService
}

func NewProfileHandler(service profile.ProfileService, homeService index.HomeService) *ProfileHandler {
	return &ProfileHandler{profileService: service, homeService: homeService}
}
