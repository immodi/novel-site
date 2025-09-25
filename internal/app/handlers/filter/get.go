package filter

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/http/templates/filter"
	"net/http"
)

func (h *FilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	meta := BuildFilterMeta()
	genres, err := h.novelService.GetAllGenres()
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	handlers.GenericHandler(w, r, meta, filter.AdvancedSearch(genres))
}
