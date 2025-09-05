package search

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *SearchHandler) SearchNovel(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	novelName := strings.ToLower(pkg.SlugToTitle(chi.URLParam(r, "novelName")))

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

	h.GenericSearch(
		"search",
		totalSearchResults,
		"search",
		novelName,
		searchNovels,
		currentPage,
		w,
		r,
	)
}
