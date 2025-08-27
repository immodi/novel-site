package novelservice

import (
	"fmt"
	"immodi/novel-site/internal/app"
	"immodi/novel-site/internal/http/routes/novels"
	"immodi/novel-site/pkg"
	"net/http"
	"strconv"
)

func NovelHandler(w http.ResponseWriter, r *http.Request) {
	// Get ?page= from query params
	pageStr := r.URL.Query().Get("page")

	// Default page = 1
	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	chapters := []novels.Chapter{}
	for i := 1; i <= 21; i++ {
		chapters = append(chapters, novels.Chapter{Title: fmt.Sprintf("Chapter %d", i), Number: i})
	}

	totalChapters := 21
	novel := novels.Novel{
		Name:                "My Novel",
		Description:         "A test novel.",
		Author:              "Author",
		Status:              "Ongoing",
		CoverImage:          "https://dummyimage.com/500x720/8a888a/ffffff",
		TotalChaptersNumber: totalChapters,
		CurrentPage:         pkg.AdjustPageNumber(currentPage, totalChapters),
		TotalPages:          pkg.CalculateTotalPages(totalChapters),
		Chapters:            pkg.GetPageChapters(chapters, currentPage),
	}

	app.GenericServiceHandler(w, r, "Test Novel", novels.NovelInfo(novel))
}
