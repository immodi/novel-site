package utils

import (
	"immodi/novel-site/internal/app/handlers/auth"
	"immodi/novel-site/internal/app/handlers/chapters"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/app/handlers/novels"
	"immodi/novel-site/internal/app/handlers/privacy"
	"immodi/novel-site/internal/app/handlers/terms"
)

type Handlers struct {
	Novel   *novels.NovelHandler
	Home    *index.HomeHandler
	Chapter *chapters.ChapterHandler
	Terms   *terms.TermsHandler
	Privacy *privacy.PrivacyHandler
	Auth    *auth.AuthHandler
}

func RegisterHandlers(svcs *Services) *Handlers {
	return &Handlers{
		Novel:   novels.NewNovelHandler(svcs.NovelService),
		Home:    index.NewHomeHandler(svcs.HomeService),
		Chapter: chapters.NewChapterHandler(svcs.ChapterService),
		Auth:    auth.NewAuthHandler(svcs.AuthService),
		Terms:   terms.NewTermsHandler(),
		Privacy: privacy.NewPrivacyHandler(),
	}
}
