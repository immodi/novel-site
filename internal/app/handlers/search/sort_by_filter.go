package search

import (
	"immodi/novel-site/internal/app/handlers"
	"log"
	"math"
	"net/http"
)

func (h *SearchHandler) SortByFilter(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handlers.NotFoundHandler(w, r)
		return
	}

	// currentPage := 1
	// if pageStr := r.Form.Get("page"); pageStr != "" {
	// 	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
	// 		currentPage = p
	// 	}
	// }

	// Genres (checkboxes – slice of strings)
	// genres := r.Form["genres"] // []string

	// Whether to combine genres with AND/OR/EXCLUDE
	// genreCondition := r.Form.Get("ctgcon") // string: "and", "or", "exclude"

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

	chapterFilteredNovels, err := h.searchService.ListNovelsByChapterRange(chapterMin, chapterMax)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	for _, novel := range chapterFilteredNovels {
		log.Println(novel.Title)
	}

	// Translation status (-1 all, 0 ongoing, 1 completed)
	// status := r.Form.Get("status")

	// Sort option
	// sortBy := r.Form.Get("sort")

	// Tags included (chips/pills)
	// tags := r.Form["tags"] // []string

	// Tags excluded (chips/pills)
	// tagsExcluded := r.Form["tagsExcluded"] // []string
	//
	// results := DbNovelsToSearchDtosMapper(h.searchService, dbResults)
	//
	// currentPage = pkg.AdjustPageNumber(currentPage, int(totalResults), pkg.SEARCH_PAGE_LIMIT)
	//
	// h.GenericSearch(
	// 	"advanced", // search type
	// 	totalResults,
	// 	"Advanced Search Results", // page title
	// 	"",                        // optional subtitle
	// 	results,
	// 	currentPage,
	// 	w,
	// 	r,
	// )
}
