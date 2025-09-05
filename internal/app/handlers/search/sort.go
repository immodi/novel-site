package search

import (
	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
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
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	dbResults, err := h.searchService.ListSortedNovels(collection, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}
	results := dbNovelsToSearchDtos(h.searchService, dbResults)

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

func (h *SearchHandler) SortNovelsByGenres(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	genreFromHeader := strings.ToLower(chi.URLParam(r, "genre"))

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	genre := sql.Genre(genreFromHeader)
	if !genre.IsValidGenre() {
		genre = sql.Genre(sql.GenreAction)
	}

	totalResults, err := h.searchService.CountNovelsByGenre(genre)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	dbResults, err := h.searchService.ListNovelsByGenre(genre, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}
	results := dbNovelsToSearchDtos(h.searchService, dbResults)

	h.GenericSearch(
		"genre",
		totalResults,
		"genre",
		strings.ToUpper(string(genre)),
		results,
		currentPage,
		w,
		r,
	)
}

func (h *SearchHandler) SortNovelsByTags(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	tag := strings.ToLower(pkg.SlugToTitle(chi.URLParam(r, "tag")))

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	totalResults, err := h.searchService.CountNovelsByTag(tag)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	dbResults, err := h.searchService.ListNovelsByTag(tag, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}
	results := dbNovelsToSearchDtos(h.searchService, dbResults)

	h.GenericSearch(
		"tag",
		totalResults,
		"tag",
		strings.ToUpper(tag),
		results,
		currentPage,
		w,
		r,
	)
}

func (h *SearchHandler) SortNovelsByAuthor(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	author := strings.ToLower(pkg.SlugToTitle(chi.URLParam(r, "author")))

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	totalResults, err := h.searchService.CountNovelsByAuthor(author)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	dbResults, err := h.searchService.ListNovelsByAuthor(author, offset, pkg.SEARCH_PAGE_LIMIT)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}
	results := dbNovelsToSearchDtos(h.searchService, dbResults)

	h.GenericSearch(
		"author",
		totalResults,
		"auhtor",
		strings.ToUpper(author),
		results,
		currentPage,
		w,
		r,
	)
}
