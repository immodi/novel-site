package search

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *SearchHandler) SortNovelsByGenres(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	genreSlug := chi.URLParam(r, "genre")

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	totalResults, err := h.searchService.CountNovelsByGenre(genreSlug)
	if err != nil {
		handlers.NotFoundHandler(w, r)
		log.Println(err.Error())
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	genreName, dbResults, err := h.searchService.ListNovelsByGenre(genreSlug, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.NotFoundHandler(w, r)
		log.Println(err.Error())
		return
	}
	results := DbNovelsToSearchDtosMapper(h.searchService, dbResults)

	h.GenericSearch(
		"genre",
		totalResults,
		"genre",
		genreName,
		results,
		currentPage,
		w,
		r,
	)
}
