package search

import (
	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *SearchHandler) SortNovelsByCollection(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	collectionFromHeader := strings.ToLower(chi.URLParam(r, "collection"))
	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	collection := sql.Collection(collectionFromHeader)
	if !collection.IsValidCollection() {
		collection = sql.Collection(sql.CollectionHot)
	}

	totalResults, err := h.searchService.CountSortedNovels(collection)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	dbResults, err := h.searchService.ListSortedNovels(collection, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}
	results := DbNovelsToSearchDtosMapper(h.searchService, dbResults)

	h.GenericSearch(
		"sort",
		totalResults,
		"sort by collection",
		strings.ToUpper(string(collection)),
		results,
		currentPage,
		w,
		r,
	)
}
