package novels

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/comments"
	"immodi/novel-site/internal/db/repositories"
	novelscomponents "immodi/novel-site/internal/http/templates/novels"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *NovelHandler) GetNovel(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	novelSlug := chi.URLParam(r, "novelSlug")
	successMsg := GetAndClearCookie(w, r, "successMessage")
	errorMsg := GetAndClearCookie(w, r, "errorMessage")

	currentPage := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			currentPage = p
		}
	}

	dbNovel, err := h.novelService.GetNovelBySlug(novelSlug)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	totalChaptersInt64, err := h.novelService.CountChapters(dbNovel.ID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	totalChapters := int(totalChaptersInt64)
	currentPage = pkg.AdjustPageNumber(currentPage, totalChapters, pkg.PAGE_LIMIT)

	// Chapters
	dbChapters, err := h.novelService.GetChapters(int(dbNovel.ID), currentPage)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	var lastChapter repositories.Chapter
	lastChapter, err = h.novelService.GetLastChapter(dbNovel.ID)
	if err != nil {
		lastChapter = repositories.Chapter{
			Title: "Chapter doesn't exist",
		}
	}

	chapters := DbChaptersToInfoChaptersMapper(dbChapters)

	genres, _ := h.novelService.GetGenres(dbNovel.ID)
	tags, _ := h.novelService.GetTags(dbNovel.ID)

	go IncrementNovelViews(h.novelService, dbNovel.ID)

	novelStatus := "Completed"
	if dbNovel.IsCompleted == 0 {
		novelStatus = "Ongoing"
	}

	isNovelBookMarked := IsNovelBookMarked(r, dbNovel.ID, h.novelService)
	novel := MapToNovel(
		dbNovel,
		genres,
		tags,
		novelStatus,
		totalChapters,
		isNovelBookMarked,
		currentPage,
		chapters,
		&lastChapter,
		successMsg,
		errorMsg,
	)

	var isRedirect bool = false
	isRedirectStr := GetAndClearCookie(w, r, comments.CommentRedirectCookie)
	if isRedirectStr != "" {
		isRedirect = true
	}

	metaData := BuildNovelMeta(novel, novelStatus)
	handlers.GenericHandler(w, r, metaData, novelscomponents.NovelInfo(novel, isRedirect))
}
