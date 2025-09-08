package index

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/http/templates/index"
	"net/http"
)

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	dbNewestNovels, err := h.homeService.ListNewestNovels()
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	dbHotNovels, err := h.homeService.ListHotNovels()
	newestNovels, err := DbNovelToHomeNovelMapper(dbNewestNovels, h.homeService)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	hotNovels, err := DbNovelToHomeNovelMapper(dbHotNovels, h.homeService)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}
	println(hotNovels[0].Slug)

	dbCompletedNovels, err := h.homeService.ListCompletedNovels()
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	// genres, err := h.homeService.ListGenres()
	// if err != nil {
	// 	handlers.ServerErrorHandler(w, r)
	// 	return
	// }

	completedNovels, err := DbNovelToHomeNovelMapper(dbCompletedNovels, h.homeService)
	handlers.GenericServiceHandler(w, r, GetIndexMetaData(), index.Index(hotNovels, newestNovels, completedNovels))
}
