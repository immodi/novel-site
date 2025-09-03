package index

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/http/templates/index"
	"net/http"
)

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	dbNewestNovels, err := h.homeService.ListNewestNovels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbHotNovels, err := h.homeService.ListHotNovels()
	newestNovels, err := dbNovelToHomeNovelMapper(dbNewestNovels, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hotNovels, err := dbNovelToHomeNovelMapper(dbHotNovels, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbCompletedNovels, err := h.homeService.ListCompletedNovels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	completedNovels, err := dbNovelToHomeNovelMapper(dbCompletedNovels, h)

	handlers.GenericServiceHandler(w, r, GetIndexMetaData(), index.Index(hotNovels, newestNovels, completedNovels))
}
