package load

import "immodi/novel-site/internal/app/services/load"

type LoadHandler struct {
	loadService load.LoadService
}

func NewLoadHandler(service load.LoadService) *LoadHandler {
	return &LoadHandler{loadService: service}
}
