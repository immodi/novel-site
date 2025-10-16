package updater

import "immodi/novel-site/internal/app/services/profile"

type UpdaterHandler struct {
	GrpcURL        string
	profileService profile.ProfileService
}

func NewUpdaterHandler(GrpcURL string, profileService profile.ProfileService) *UpdaterHandler {
	return &UpdaterHandler{GrpcURL: GrpcURL, profileService: profileService}
}
