package handlers

import (
	handlers "immodi/novel-site/internal/app/handlers/paths"
	"immodi/novel-site/internal/app/services"
)

type Handlers struct {
	Novel   *handlers.NovelHandler
	Home    *handlers.HomeHandler
	Chapter *handlers.ChapterHandler
}

func RegisterHandlers(svcs *services.Services) *Handlers {
	return &Handlers{
		Novel:   handlers.NewNovelHandler(svcs.DB),
		Home:    handlers.NewHomeHandler(svcs.DB),
		Chapter: handlers.NewChapterHandler(svcs.DB),
	}
}
