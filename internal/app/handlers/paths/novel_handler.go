package handlers

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services"
	repositories "immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/templates/components"
	"immodi/novel-site/internal/http/templates/novels"
	"immodi/novel-site/pkg"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type NovelHandler struct {
	dbService *services.DBService
}

func NewNovelHandler(dbService *services.DBService) *NovelHandler {
	return &NovelHandler{dbService: dbService}
}

func (h *NovelHandler) GetNovel(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	novelName := pkg.SlugToTitle(chi.URLParam(r, "novelName"))

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	dbNovel, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByNameLike(ctx, fmt.Sprintf("%s", novelName)) // exact match
		// return q.GetNovelByNameLike(ctx, fmt.Sprintf("%%%s%%", novelName)) // partial match
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbChapters, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]repositories.Chapter, error) {
		return q.ListChaptersByNovelPaginated(ctx, repositories.ListChaptersByNovelPaginatedParams{
			NovelID: dbNovel.ID,
			Limit:   pkg.PAGE_LIMIT,
			Offset:  int64(pkg.PAGE_LIMIT * (currentPage - 1)),
		})
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalChaptersInt64, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountChaptersByNovel(ctx, dbNovel.ID)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalChapters := int(totalChaptersInt64)
	chapters := castDbChaptersToInfoChapters(dbChapters)

	genres, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		genreRows, err := q.ListGenresByNovel(ctx, 1) // novel ID
		if err != nil {
			return nil, err
		}

		// Extract just the genre strings from the rows
		var genres []string
		for _, row := range genreRows {
			genres = append(genres, row)
		}
		return genres, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	novel := novels.Novel{
		Name:                dbNovel.Title,
		Description:         dbNovel.Description,
		Author:              dbNovel.Author,
		Genres:              genres,
		Status:              dbNovel.Status,
		CoverImage:          dbNovel.CoverImage,
		TotalChaptersNumber: totalChapters,
		CurrentPage:         pkg.AdjustPageNumber(currentPage, totalChapters),
		TotalPages:          pkg.CalculateTotalPages(totalChapters),
		Chapters:            chapters,
		LastChapterName:     dbNovel.LatestChapterName,
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
		UpdateTime:        dbNovel.UpdateTime,
		LatestChapterName: novel.LastChapterName,
		LatestChapterURL:  fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, novel.Name, novel.TotalChaptersNumber),
	}

	GenericServiceHandler(w, r, metaData, novels.NovelInfo(novel))
}

func (h *NovelHandler) CreateNovel(w http.ResponseWriter, r *http.Request) {
	// Create a novel with default values
	dbNovel, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.CreateNovel(ctx, repositories.CreateNovelParams{
			Title:               "Test Novel",
			Description:         "This is a default novel description. An exciting story awaits!",
			CoverImage:          "https://dummyimage.com/500x720/8a818a/2fffff",
			Author:              "Default Author 2",
			Status:              "Ongoing",
			UpdateTime:          "2025-08-30",
			LatestChapterName:   "Chapter 1",
			TotalChaptersNumber: 1,
		})
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create novel: %v", err), http.StatusInternalServerError)
		return
	}

	// Add default genres
	defaultGenres := []string{"Adventure", "Drama", "Fantasy"}
	for _, genre := range defaultGenres {
		_, err = services.ExecuteWithResult[any](h.dbService, func(ctx context.Context, q *repositories.Queries) (any, error) {
			return nil, q.AddGenreToNovel(ctx, repositories.AddGenreToNovelParams{
				NovelID: dbNovel.ID,
				Genre:   genre,
			})
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add genre %s: %v", genre, err), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/novel", http.StatusSeeOther)
}

func castDbChaptersToInfoChapters(dbChapters []repositories.Chapter) []novels.Chapter {
	var chapters []novels.Chapter
	for _, dbChapter := range dbChapters {
		chapters = append(chapters, novels.Chapter{
			Title: dbChapter.Title,
		})
	}
	return chapters
}
