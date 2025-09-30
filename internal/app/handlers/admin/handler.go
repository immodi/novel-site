package admin

import (
	"immodi/novel-site/internal/app/services/auth"
	"immodi/novel-site/internal/app/services/chapters"
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/app/services/profile"
)

type AdminHandler struct {
	authService    auth.AuthService
	indexService   index.HomeService
	chapterService chapters.ChapterService
	profileService profile.ProfileService
}

func NewAdminHandler(
	service auth.AuthService,
	profileService profile.ProfileService,
	indexService index.HomeService,
	chapterService chapters.ChapterService,
) *AdminHandler {
	return &AdminHandler{
		authService:    service,
		profileService: profileService,
		indexService:   indexService,
		chapterService: chapterService,
	}
}
