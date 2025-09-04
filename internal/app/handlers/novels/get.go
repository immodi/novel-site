package novels

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/db/repositories"
	novelscomponents "immodi/novel-site/internal/http/templates/novels"
	"immodi/novel-site/pkg"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *NovelHandler) GetNovel(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	novelName := pkg.SlugToTitle(chi.URLParam(r, "novelName"))

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	dbNovel, err := h.novelService.GetNovelByName(novelName)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	totalChaptersInt64, err := h.novelService.CountChapters(dbNovel.ID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	totalChapters := int(totalChaptersInt64)
	currentPage = pkg.AdjustPageNumber(currentPage, totalChapters, pkg.PAGE_LIMIT)

	// Chapters
	dbChapters, err := h.novelService.GetChapters(int(dbNovel.ID), currentPage)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	var lastChapter repositories.Chapter
	lastChapter, err = h.novelService.GetLastChapter(dbNovel.ID)
	if err != nil {
		lastChapter = repositories.Chapter{
			Title: "Chapter doesn't exist",
		}
	}

	chapters := CastDbChaptersToInfoChapters(dbChapters)

	genres, _ := h.novelService.GetGenres(dbNovel.ID)
	tags, _ := h.novelService.GetTags(dbNovel.ID)

	go IncrementNovelViews(h.novelService, dbNovel.ID)

	novelStatus := "Completed"
	if dbNovel.IsCompleted == 0 {
		novelStatus = "Ongoing"
	}

	novel := MapDBNovelToNovel(dbNovel, genres, tags, novelStatus, totalChapters, currentPage, chapters, &lastChapter)
	metaData := MapNovelToMetaData(*novel, novelStatus)

	handlers.GenericServiceHandler(w, r, metaData, novelscomponents.NovelInfo(*novel))
}
