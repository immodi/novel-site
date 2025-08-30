package handlers

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services"
	repositories "immodi/novel-site/internal/db/repositories"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ChapterHandler struct {
	dbService *services.DBService
}

func NewChapterHandler(dbService *services.DBService) *ChapterHandler {
	return &ChapterHandler{dbService: dbService}
}

// ReadChapter handles GET /novel/{novel_name}/chapter-{number}
// func (h *ChapterHandler) ReadChapter(w http.ResponseWriter, r *http.Request) {
// 	// Extract chapter number from URL (you'll need to adapt this to your routing)
// 	chapterNumStr := r.URL.Query().Get("chapter")
// 	if chapterNumStr == "" {
// 		http.Error(w, "Chapter number required", http.StatusBadRequest)
// 		return
// 	}
//
// 	chapterNum, err := strconv.Atoi(chapterNumStr)
// 	if err != nil || chapterNum < 1 {
// 		http.Error(w, "Invalid chapter number", http.StatusBadRequest)
// 		return
// 	}
//
// 	novelID := int64(1) // You'll need to get this from URL or context
//
// 	// Get the novel first
// 	dbNovel, err := services.ExecuteWithResult[repositories.Novel](h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
// 		return q.GetNovelByID(ctx, novelID)
// 	})
// 	if err != nil {
// 		http.Error(w, "Novel not found", http.StatusNotFound)
// 		return
// 	}
//
// 	// Get the chapter
// 	dbChapter, err := services.ExecuteWithResult[repositories.Chapter](h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
// 		return q.GetChapterByNovelAndNumber(ctx, repositories.GetChapterByNovelAndNumberParams{
// 			NovelID: novelID,
// 			Number:  int64(chapterNum),
// 		})
// 	})
// 	if err != nil {
// 		http.Error(w, "Chapter not found", http.StatusNotFound)
// 		return
// 	}
//
// 	// Convert to template struct
// 	chapter := chapters.Chapter{
// 		Title:       dbChapter.Title,
// 		Number:      int(dbChapter.Number),
// 		Content:     services.GetStringValue(dbChapter.Content),
// 		NovelName:   dbNovel.Title,
// 		NovelAuthor: services.GetStringValue(dbNovel.Author),
// 	}
//
// 	// Create metadata
// 	metaData := &components.MetaDataStruct{
// 		IsRendering: true,
// 		Title:       fmt.Sprintf("%s - Chapter %d - %s", dbNovel.Title, chapterNum, components.SITE_NAME),
// 		Description: fmt.Sprintf("Read Chapter %d of %s by %s", chapterNum, dbNovel.Title, services.GetStringValue(dbNovel.Author)),
// 		Keywords:    fmt.Sprintf("%s chapter %d, read %s chapter %d online", dbNovel.Title, chapterNum, dbNovel.Title, chapterNum),
// 		OgURL:       fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, dbNovel.Title, chapterNum),
// 		Canonical:   fmt.Sprintf("%s/novel/%s/chapter-%d", components.DOMAIN, dbNovel.Title, chapterNum),
// 		NovelName:   dbNovel.Title,
// 	}
//
// 	GenericServiceHandler(w, r, metaData, chapters.ChapterRead(chapter))
// }

// CreateChapter handles POST /admin/chapter/create
// func (h *ChapterHandler) CreateChapter(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	// Parse form data
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Failed to parse form", http.StatusBadRequest)
// 		return
// 	}
//
// 	// Get form values
// 	novelIDStr := r.FormValue("novel_id")
// 	title := r.FormValue("title")
// 	content := r.FormValue("content")
// 	numberStr := r.FormValue("number")
//
// 	// Validate inputs
// 	if novelIDStr == "" || title == "" || numberStr == "" {
// 		http.Error(w, "Missing required fields: novel_id, title, number", http.StatusBadRequest)
// 		return
// 	}
//
// 	novelID, err := strconv.ParseInt(novelIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid novel ID", http.StatusBadRequest)
// 		return
// 	}
//
// 	number, err := strconv.ParseInt(numberStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid chapter number", http.StatusBadRequest)
// 		return
// 	}
//
// 	// Verify novel exists
// 	_, err = services.ExecuteWithResult[repositories.Novel](h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
// 		return q.GetNovelByID(ctx, novelID)
// 	})
// 	if err != nil {
// 		http.Error(w, "Novel not found", http.StatusNotFound)
// 		return
// 	}
//
// 	// Create the chapter
// 	dbChapter, err := services.ExecuteWithResult[repositories.Chapter](h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
// 		return q.CreateChapter(ctx, repositories.CreateChapterParams{
// 			NovelID: novelID,
// 			Title:   title,
// 			Content: content,
// 		})
// 	})
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to create chapter: %v", err), http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Update novel's latest chapter info
// 	_, err = services.ExecuteWithResult[repositories.Novel](h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
// 		return q.UpdateNovel(ctx, repositories.UpdateNovelParams{
// 			ID:                  novelID,
// 			Title:               "", // You'll need to get current values first
// 			Description:         services.NewNullString(""),
// 			CoverImage:          services.NewNullString(""),
// 			Author:              services.NewNullString(""),
// 			Status:              services.NewNullString(""),
// 			UpdateTime:          services.NewNullString("2025-08-30"),
// 			LatestChapterName:   services.NewNullString(title),
// 			TotalChaptersNumber: services.NewNullInt64(number),
// 		})
// 	})
//
// 	// Redirect to the created chapter or return success
// 	http.Redirect(w, r, fmt.Sprintf("/novel/%d/chapter-%d", novelID, number), http.StatusSeeOther)
// }

// CreateChapterWithDefaults creates a chapter with default content
func (h *ChapterHandler) CreateChapterWithDefaults(w http.ResponseWriter, r *http.Request) {
	novelIDString := chi.URLParam(r, "novelId")
	novelID, err := strconv.ParseInt(novelIDString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid novel ID", http.StatusBadRequest)
		return
	}

	// Get current chapter count to determine next chapter number
	chapters, err := services.ExecuteWithResult[[]repositories.Chapter](h.dbService, func(ctx context.Context, q *repositories.Queries) ([]repositories.Chapter, error) {
		return q.ListChaptersByNovel(ctx, novelID)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get chapters: %v", err), http.StatusInternalServerError)
		return
	}

	nextChapterNum := int64(len(chapters) + 1)

	// Create chapter with defaults
	_, err = services.ExecuteWithResult(h.dbService, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.CreateChapter(ctx, repositories.CreateChapterParams{
			NovelID: novelID,
			Title:   fmt.Sprintf("Chapter %d", nextChapterNum),
			Content: "test chapter 1",
		})
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create chapter: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	http.Redirect(w, r, "/novel", http.StatusSeeOther)
}
