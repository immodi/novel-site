package chapters

import (
	"immodi/novel-site/internal/app/services/chapters"
)

type ChapterHandler struct {
	chapterService chapters.ChapterService
}

func NewChapterHandler(service chapters.ChapterService) *ChapterHandler {
	return &ChapterHandler{chapterService: service}
}
