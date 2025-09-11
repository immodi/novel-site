package utils

import (
	"immodi/novel-site/internal/app/handlers/auth"
	"immodi/novel-site/internal/app/handlers/chapters"
	"immodi/novel-site/internal/app/handlers/comments"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/app/handlers/load"
	"immodi/novel-site/internal/app/handlers/novels"
	"immodi/novel-site/internal/app/handlers/privacy"
	"immodi/novel-site/internal/app/handlers/profile"
	"immodi/novel-site/internal/app/handlers/search"
	"immodi/novel-site/internal/app/handlers/terms"
)

type Handlers struct {
	Novel   *novels.NovelHandler
	Home    *index.HomeHandler
	Chapter *chapters.ChapterHandler
	Terms   *terms.TermsHandler
	Privacy *privacy.PrivacyHandler
	Auth    *auth.AuthHandler
	Search  *search.SearchHandler
	Load    *load.LoadHandler
	Profile *profile.ProfileHandler
	Comment *comments.CommentHandler
}

func RegisterHandlers(svcs *Services) *Handlers {
	return &Handlers{
		Novel:   novels.NewNovelHandler(svcs.NovelService),
		Home:    index.NewHomeHandler(svcs.HomeService),
		Chapter: chapters.NewChapterHandler(svcs.ChapterService),
		Search:  search.NewSearchHandler(svcs.SearchServie, svcs.HomeService),
		Auth:    auth.NewAuthHandler(svcs.AuthService),
		Terms:   terms.NewTermsHandler(),
		Privacy: privacy.NewPrivacyHandler(),
		Load:    load.NewLoadHandler(svcs.LoadService),
		Profile: profile.NewProfileHandler(svcs.ProfileService, svcs.HomeService),
		Comment: comments.NewCommentHandler(svcs.CommentService, svcs.ProfileService),
	}
}
