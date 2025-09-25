package filter

import "immodi/novel-site/internal/app/services/novels"

type FilterHandler struct {
	novelService novels.NovelService
}

func NewFilterHandler(service novels.NovelService) *FilterHandler {
	return &FilterHandler{novelService: service}
}
