package pkg

import "fmt"

const PAGE_LIMIT = 21
const SEARCH_PAGE_LIMIT = 8
const PROFILE_PAGE_LIMIT = 6

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

func MakePagesCompact(current, total int) []string {
	const (
		maxSmallPages = 2 // if total pages <= this, show all pages
		windowSize    = 1 // how many pages to show on each side of current
	)

	// Small case: show everything
	if total <= maxSmallPages {
		pages := []string{}
		for i := 1; i <= total; i++ {
			pages = append(pages, fmt.Sprintf("%d", i))
		}
		return pages
	}

	pages := []string{}

	// Sliding window around current
	start := max(current-windowSize, 1)
	end := min(current+windowSize, total)

	// Add the window pages
	for i := start; i <= end; i++ {
		pages = append(pages, fmt.Sprintf("%d", i))
	}

	// If there are still pages left after the window, add ellipsis and last page
	if end < total {
		if end < total-1 {
			pages = append(pages, "…")
		}
		pages = append(pages, fmt.Sprintf("%d", total))
	}

	return pages
}
