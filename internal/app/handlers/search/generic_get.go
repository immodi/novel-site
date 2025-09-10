package search

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/index"
	sql "immodi/novel-site/internal/db/schema"
	searchresutlsdto "immodi/novel-site/internal/http/structs/search"
	"immodi/novel-site/internal/http/templates/search"
	"immodi/novel-site/pkg"
	"math"
	"net/http"
)

func (h *SearchHandler) GenericSearch(paramter string, totalResults int64, headerText string, query string, collection []searchresutlsdto.SearchResultDto, currentPage int, w http.ResponseWriter, r *http.Request) {
	dbHotNovels, err := h.searchService.ListSortedNovels(sql.CollectionHot, 0, 5)
	hotNovels, err := index.DbNovelToHomeNovelMapper(dbHotNovels, h.homeService)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	totalPages := int(math.Ceil(float64(totalResults) / float64(pkg.SEARCH_PAGE_LIMIT)))

	handlers.GenericHandler(w, r, index.BuildHomeMeta(), search.SearchResults(paramter, headerText, query, collection, hotNovels, currentPage, totalPages))
}
