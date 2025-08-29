package novelservice

import (
	"fmt"
	"immodi/novel-site/internal/app"
	"immodi/novel-site/internal/http/templates/components"
	"immodi/novel-site/internal/http/templates/novels"
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
		Genres:              []string{"Adventure", "Drama", "Fantasy"},
		TotalChaptersNumber: totalChapters,
		CurrentPage:         pkg.AdjustPageNumber(currentPage, totalChapters),
		TotalPages:          pkg.CalculateTotalPages(totalChapters),
		Chapters:            pkg.GetPageChapters(chapters, currentPage),
		LastChapterName:     "Chapter 21",
	}
	metaData := &components.MetaDataStruct{
		IsRendering:       true,
		Title:             fmt.Sprintf("%s - Read %s For Free - %s", novel.Name, novel.Name, components.SITE_NAME),
		Description:       novel.Description,
		Keywords:          fmt.Sprintf("%s novel 2025, read %s online 2025, free %s novel", novel.Name, novel.Name, novel.Name),
		OgURL:             fmt.Sprintf("%s/novel/%s", components.DOMAIN, novel.Name),
		Canonical:         fmt.Sprintf("%s/novel/%s", components.DOMAIN, novel.Name),
		CoverImage:        novel.CoverImage,
		Genres:            novel.Genres,
		Author:            novel.Author,
		Status:            novel.Status,
		AuthorLink:        fmt.Sprintf("%s/author/%s", components.DOMAIN, novel.Author),
		NovelName:         novel.Name,
		ReadURL:           fmt.Sprintf("%s/novel/%s/chapter-1", components.DOMAIN, novel.Name),
		UpdateTime:        "2025-08-29T12:00:00Z",
		LatestChapterName: novel.LastChapterName,
		LatestChapterURL:  fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, novel.Name, novel.TotalChaptersNumber),
	}

	app.GenericServiceHandler(w, r, metaData, novels.NovelInfo(novel))
}
