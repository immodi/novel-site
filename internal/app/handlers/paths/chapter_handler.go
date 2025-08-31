package handlers

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services"
	repositories "immodi/novel-site/internal/db/repositories"
	novels "immodi/novel-site/internal/http/templates/chapters"
	"immodi/novel-site/internal/http/templates/components"
	cnovels "immodi/novel-site/internal/http/templates/novels"
	"immodi/novel-site/pkg"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type ChapterHandler struct {
	dbService *services.DBService
}

func NewChapterHandler(dbService *services.DBService) *ChapterHandler {
	return &ChapterHandler{dbService: dbService}
}

func (h *ChapterHandler) ReadChapter(w http.ResponseWriter, r *http.Request) {
	novelName := pkg.SlugToTitle(chi.URLParam(r, "novelName"))
	chapterNumStr := pkg.SlugToTitle(chi.URLParam(r, "chapterNumber"))

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil || chapterNum < 1 {
		http.Error(w, "Invalid chapter number", http.StatusBadRequest)
		return
	}

	dbNovel, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByNameLike(ctx, novelName)
	})
	if err != nil {
		http.Error(w, "Novel not found", http.StatusNotFound)
		return
	}

	totalChaptersNumber, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountChaptersByNovel(ctx, dbNovel.ID)
	})

	if err != nil {
		http.Error(w, "Failed to get total chapters", http.StatusInternalServerError)
		return
	}

	// Get the chapter
	dbChapter, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.GetChapterByNumber(ctx, repositories.GetChapterByNumberParams{
			NovelID:       dbNovel.ID,
			ChapterNumber: int64(chapterNum),
		})
	})
	if err != nil {
		http.Error(w, "Chapter not found", http.StatusNotFound)
		return
	}

	var nextChapterIntPointer *int
	var prevChapterIntPointer *int

	if dbChapter.ChapterNumber+1 > totalChaptersNumber {
		nextChapterIntPointer = nil
	} else {
		nextChapterInt := int(dbChapter.ChapterNumber + 1)
		nextChapterIntPointer = &nextChapterInt
	}

	if dbChapter.ChapterNumber-1 < 1 {
		prevChapterIntPointer = nil
	} else {
		prevChapterInt := int(dbChapter.ChapterNumber - 1)
		prevChapterIntPointer = &prevChapterInt
	}

	// Convert to template struct
	chapter := novels.ChapterPage{
		NovelName:      dbNovel.Title,
		ChapterTitle:   dbChapter.Title,
		ChapterContent: dbChapter.Content,
		PrevChapter:    prevChapterIntPointer,
		NextChapter:    nextChapterIntPointer,
	}

	// Create chapter-specific metadata
	metaData := &components.MetaDataStruct{
		IsRendering: true,
		Title:       fmt.Sprintf("%s - Chapter %d - Read %s Online - %s", dbNovel.Title, chapterNum, dbNovel.Title, components.SITE_NAME),
		Description: fmt.Sprintf("Read %s Chapter %d online. %s - Chapter %d by %s for free in high quality.", dbNovel.Title, chapterNum, dbNovel.Title, chapterNum, dbNovel.Author),
		Keywords:    fmt.Sprintf("read %s chapter %d online, free %s chapter %d, %s novel chapter %d", dbNovel.Title, chapterNum, dbNovel.Title, chapterNum, dbNovel.Title, chapterNum),
		OgURL:       fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, pkg.TitleToSlug(dbNovel.Title), chapterNum),
		Canonical:   fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, pkg.TitleToSlug(dbNovel.Title), chapterNum),

		// Extra data for meta and structured output
		CoverImage: fmt.Sprintf("%s/media/novel/%s.jpg", components.DOMAIN, pkg.TitleToSlug(dbNovel.Title)),
		Author:     dbNovel.Author,
		Status:     dbNovel.Status,

		AuthorLink: fmt.Sprintf("%s/author/%s", components.DOMAIN, pkg.TitleToSlug(dbNovel.Author)),
		NovelName:  dbNovel.Title,

		// Navigation metadata
		ReadURL:    fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, pkg.TitleToSlug(dbNovel.Title), chapterNum),
		UpdateTime: dbNovel.UpdateTime,
		// LatestChapterName: dbNovel.LastChapterName,
		// LatestChapterURL:  fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, pkg.TitleToSlug(dbNovel.Title), dbNovel.TotalChaptersNumber),
	}

	GenericServiceHandler(w, r, metaData, novels.ChapterReader(chapter))
}

// CreateChapterWithDefaults creates a chapter with default content
func (h *ChapterHandler) CreateChapterWithDefaults(w http.ResponseWriter, r *http.Request) {
	novelIDString := chi.URLParam(r, "novelId")
	novelID, err := strconv.ParseInt(novelIDString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid novel ID", http.StatusBadRequest)
		return
	}

	// Get current chapter count to determine next chapter number
	currentChaptersCount, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountChaptersByNovel(ctx, novelID)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get chapters: %v", err), http.StatusInternalServerError)
		return
	}

	// Create chapter with defaults
	_, err = services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.CreateChapter(ctx, repositories.CreateChapterParams{
			NovelID:       novelID,
			Title:         fmt.Sprintf("Chapter %d", currentChaptersCount+1),
			ChapterNumber: currentChaptersCount + 1,
			Content:       pkg.LoremText(40),
		})
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create chapter: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	http.Redirect(w, r, "/novel", http.StatusSeeOther)
}

func (h *ChapterHandler) GetChaptersDropDown(w http.ResponseWriter, r *http.Request) {
	novelName := pkg.SlugToTitle(chi.URLParam(r, "novelName"))

	dbNovel, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByNameLike(ctx, novelName)
	})
	if err != nil {
		http.Error(w, "Novel not found", http.StatusNotFound)
		return
	}

	dbChapters, err := services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) ([]repositories.Chapter, error) {
		return q.ListChaptersByNovel(ctx, dbNovel.ID)
	})

	if err != nil {
		http.Error(w, "Failed to get total chapters", http.StatusInternalServerError)
		return
	}

	chapters := dbChaptersToChapters(dbChapters)

	cmp := novels.ChapterDropdown(pkg.TitleToSlug(novelName), chapters)
	templ.Handler(cmp).ServeHTTP(w, r)

}

func dbChaptersToChapters(dbChapters []repositories.Chapter) []cnovels.Chapter {
	chapters := make([]cnovels.Chapter, len(dbChapters))

	for i, dbChapter := range dbChapters {
		chapters[i] = cnovels.Chapter{
			Title:  dbChapter.Title,
			Number: int(dbChapter.ChapterNumber),
		}
	}

	return chapters
}
