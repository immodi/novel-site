package index

import (
	"immodi/novel-site/internal/app/services/index"
)

type HomeHandler struct {
	homeService index.HomeService
}

func NewHomeHandler(service index.HomeService) *HomeHandler {
	return &HomeHandler{homeService: service}
}
