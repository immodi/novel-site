package handlers

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services"
	repositories "immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/templates/components"
	novelscomponents "immodi/novel-site/internal/http/templates/novels"
	novels "immodi/novel-site/internal/http/templates/novels/components"
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

	lastChapter, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.GetChapterByNumber(ctx, repositories.GetChapterByNumberParams{
			NovelID:       dbNovel.ID,
			ChapterNumber: totalChaptersInt64,
		})
	})

	totalChapters := int(totalChaptersInt64)
	chapters := castDbChaptersToInfoChapters(dbChapters)

	genres, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		genreRows, err := q.ListGenresByNovel(ctx, dbNovel.ID) // novel ID
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

	tags, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		tagRows, err := q.ListTagsByNovel(ctx, dbNovel.ID) // novel ID
		if err != nil {
			return nil, err
		}

		// Extract just the tag strings from the rows
		var tags []string
		for _, row := range tagRows {
			tags = append(tags, row)
		}
		return tags, nil
	})

	go incrementNovelViews(dbNovel.ID, h.dbService)

	novel := novels.Novel{
		Name:                dbNovel.Title,
		Description:         dbNovel.Description,
		Author:              dbNovel.Author,
		Genres:              genres,
		Tags:                tags,
		Status:              dbNovel.Status,
		CoverImage:          dbNovel.CoverImage,
		TotalChaptersNumber: totalChapters,
		CurrentPage:         pkg.AdjustPageNumber(currentPage, totalChapters),
		TotalPages:          pkg.CalculateTotalPages(totalChapters),
		Chapters:            chapters,
		LastChapterName:     lastChapter.Title,
		LastUpdated:         dbNovel.UpdateTime,
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

	GenericServiceHandler(w, r, metaData, novelscomponents.NovelInfo(novel))
}

func incrementNovelViews(novelId int64, dbService *services.DBService) {
	err := services.Execute(dbService, func(ctx context.Context, q *repositories.Queries) error {
		return q.IncrementNovelViewCount(ctx, novelId)
	})
	if err != nil {
		fmt.Printf("Failed to increment novel views: %v\n", err)
	}
}

func (h *NovelHandler) CreateNovelWithDefaults(w http.ResponseWriter, r *http.Request) {
	novelName := chi.URLParam(r, "novelName")

	dbNovel, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.CreateNovel(ctx, repositories.CreateNovelParams{
			Title:       novelName,
			Description: "This is a default novel description.",
			CoverImage:  "https://dummyimage.com/500x720/8a818a/ffffff",
			Author:      "Default Author 1",
			Status:      "Ongoing",
			UpdateTime:  pkg.GetCurrentTimeRFC3339(),
		})
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create novel: %v", err), http.StatusInternalServerError)
		return
	}

	// Add default genres
	defaultGenres := []string{"Adventure", "Drama", "Fantasy"}
	for _, genre := range defaultGenres {
		_, err = services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (any, error) {
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

	defaultTags := []string{"Harem", "Male Protagonist", "Magic", "Dragons"}

	for _, tag := range defaultTags {
		_, err = services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (any, error) {
			return nil, q.AddTagToNovel(ctx, repositories.AddTagToNovelParams{
				NovelID: dbNovel.ID,
				Tag:     tag,
			})
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add tag %s: %v", tag, err), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/novel", http.StatusSeeOther)
}

func castDbChaptersToInfoChapters(dbChapters []repositories.Chapter) []novels.Chapter {
	var chapters []novels.Chapter
	for _, dbChapter := range dbChapters {
		chapters = append(chapters, novels.Chapter{
			Title:  dbChapter.Title,
			Number: int(dbChapter.ChapterNumber),
		})
	}
	return chapters
}
