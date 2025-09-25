package search

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *SearchHandler) SortNovelsByTags(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	tag := chi.URLParam(r, "tag")

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	totalResults, err := h.searchService.CountNovelsByTag(tag)
	if err != nil {
		handlers.NotFoundHandler(w, r)
		log.Println(err.Error())
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	tagName, dbResults, err := h.searchService.ListNovelsByTag(tag, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.NotFoundHandler(w, r)
		log.Println(err.Error())
		return
	}
	results := DbNovelsToSearchDtosMapper(h.searchService, dbResults)

	h.GenericSearch(
		"tag",
		totalResults,
		"tag",
		tagName,
		"",
		results,
		currentPage,
		w,
		r,
	)
}
