package pkg

const PAGE_LIMIT = 21
const SEARCH_PAGE_LIMIT = 8

func MakePages(total int) []int {
	pages := make([]int, total)
	for i := 1; i <= total; i++ {
		pages[i-1] = i
	}
	return pages
}

// Calculate total pages from total chapters
func CalculateTotalPages(totalChapters int) int {
	if totalChapters == 0 {
		return 1
	}
	pages := totalChapters / PAGE_LIMIT
	if totalChapters%PAGE_LIMIT != 0 {
		pages++
	}
	return pages
}

// Return the chapters belonging to the given page
func GetPageChapters[T any](items []T, page int) []T {
	if page < 1 {
		page = 1
	}

	totalPages := (len(items) + PAGE_LIMIT - 1) / PAGE_LIMIT
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * PAGE_LIMIT
	end := min(start+PAGE_LIMIT, len(items))
	return items[start:end]
}

// AdjustPage ensures the requested page is within valid bounds.
func AdjustPageNumber(requestedPage, totalChapters, pageLimit int) int {
	if requestedPage < 1 {
		return 1
	}

	totalPages := (totalChapters + pageLimit - 1) / pageLimit
	if totalPages == 0 {
		return 1 // no chapters, but default to page 1
	}

	if requestedPage > totalPages {
		return totalPages
	}

	return requestedPage
}
