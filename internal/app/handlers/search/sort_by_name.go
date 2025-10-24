package search

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strconv"
)

func (h *SearchHandler) SortNovelsByNames(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	totalResults, err := h.searchService.CountAllNovels()
	if err != nil {
		handlers.NotFoundHandler(w, r)
		log.Println(err.Error())
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	dbResults, err := h.searchService.ListNovelsByName(offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.NotFoundHandler(w, r)
		log.Println(err.Error())
		return
	}
	results := DbNovelsToSearchDtosMapper(h.searchService, dbResults)

	h.GenericSearch(
		"",
		totalResults,
		"",
		"All Novels",
		"",
		results,
		currentPage,
		w,
		r,
	)
}
