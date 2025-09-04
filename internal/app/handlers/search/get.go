package search

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/http/templates/search"
	"immodi/novel-site/pkg"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *SearchHandler) SearchNovel(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	novelName := pkg.SlugToTitle(chi.URLParam(r, "novelName"))
	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	totalSearchResults, err := h.searchService.CountTotalSearchedNovels(novelName)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalSearchResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	dbSearchResults, err := h.searchService.SearchNovelsPaginated(novelName, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}
	searchNovels := dbNovelsToSearchDtos(h.searchService, dbSearchResults)

	dbHotNovels, err := h.homeService.ListHotNovels()
	hotNovels, err := index.DbNovelToHomeNovelMapper(dbHotNovels, h.homeService)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	totalPages := int(math.Ceil(float64(totalSearchResults) / float64(pkg.SEARCH_PAGE_LIMIT)))

	handlers.GenericServiceHandler(w, r, index.GetIndexMetaData(), search.SearchResults(novelName, searchNovels, hotNovels, currentPage, totalPages))
}
