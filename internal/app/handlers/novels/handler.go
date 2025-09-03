package novels

import (
	"immodi/novel-site/internal/app/services/novels"
)

type NovelHandler struct {
	novelService novels.NovelService
}

func NewNovelHandler(service novels.NovelService) *NovelHandler {
	return &NovelHandler{novelService: service}
}
