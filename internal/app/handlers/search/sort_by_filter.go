package search

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/pkg"
	"log"
	"math"
	"net/http"
	"strconv"
)

func (h *SearchHandler) SortByFilter(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handlers.NotFoundHandler(w, r)
		return
	}

	currentPage := 1
	if pageStr := r.Form.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}
	offset := (currentPage - 1) * pkg.SEARCH_PAGE_LIMIT

	// Genres (checkboxes – slice of strings)
	genres := r.Form["genres"] // []string

	// Whether to combine genres with AND/OR/EXCLUDE
	genreCondition := r.Form.Get("ctgcon") // string: "and", "or", "exclude"

	// Total chapter range
	totalChapter := r.Form.Get("totalchapter") // e.g. "100,200" or "0"
	chapterMin, chapterMax := 0, math.MaxInt
	minPtr, maxPtr, err := ParseTotalChapterRange(totalChapter)
	if err != nil {
		handlers.NotFoundHandler(w, r)
		log.Println("parse range:", err)
		return
	}
	if minPtr != nil {
		chapterMin = *minPtr
	}
	if maxPtr != nil {
		chapterMax = *maxPtr
	}

	// Translation status (-1 all, 0 ongoing, 1 completed)
	status := r.Form.Get("status")

	// Sort option
	sortBy := r.Form.Get("sort")

	tagCondition := r.Form.Get("tagcon") // string: "and", "or"

	// Tags included (chips/pills)
	tags := r.Form["tags"] // []string

	// Tags excluded (chips/pills)
	tagsExcluded := r.Form["tagsExcluded"] // []string
	novelsResult, err := h.searchService.FilterNovels(repositories.FilterNovelsParams{
		MinChapters:    int64(chapterMin),
		MaxChapters:    int64(chapterMax),
		TagCondition:   tagCondition,
		Tags:           tags,
		TagsExclude:    tagsExcluded,
		Genres:         genres,
		GenreCondition: genreCondition,
		IsCompleted:    int(ParseNovelStatus(status)),
		SortBy:         sortBy,
		Limit:          int64(pkg.SEARCH_PAGE_LIMIT),
		Offset:         int64(offset),
	})
	if err != nil {
		log.Println(err.Error())
		handlers.ServerErrorHandler(w, r)
		return
	}

	totalResults := novelsResult.TotalCount
	results := DbFilteredNovelsToSearchDtosMapper(h.searchService, novelsResult.Novels)
	currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	h.GenericSearch(
		"filter",
		int64(totalResults),
		"Advanced Search Results",
		"",
		baseURL(r),
		results,
		currentPage,
		w,
		r,
	)
}

func baseURL(r *http.Request) string {
	// Make a copy of the URL
	url := *r.URL
	// Clear the "page" query parameter
	q := url.Query()
	q.Del("page")
	url.RawQuery = q.Encode()
	return url.String()
}
