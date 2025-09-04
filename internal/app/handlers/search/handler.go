package search

import (
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/app/services/search"
)

type SearchHandler struct {
	searchService search.SearchService
	homeService   index.HomeService
}

func NewSearchHandler(searchService search.SearchService, homeService index.HomeService) *SearchHandler {
	return &SearchHandler{searchService: searchService, homeService: homeService}
}
