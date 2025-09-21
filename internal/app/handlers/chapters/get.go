package chapters

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/comments"
	dropdown "immodi/novel-site/internal/http/components/chapters"
	"immodi/novel-site/internal/http/templates/chapters"
	"immodi/novel-site/pkg"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func (h *ChapterHandler) ReadChapter(w http.ResponseWriter, r *http.Request) {
	novelSlug := chi.URLParam(r, "novelSlug")
	chapterNumStr := chi.URLParam(r, "chapterNumber")
	userID := pkg.IsAuthedUser(r)

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil || chapterNum < 1 {
		handlers.ServerErrorHandler(w, r)
		return
	}

	dbNovel, err := h.chapterService.GetNovelBySlug(novelSlug)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	go h.chapterService.IncrementNovelViewCount(dbNovel.ID)

	totalChaptersNumber, err := h.chapterService.CountChaptersByNovel(dbNovel.ID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	// Get the chapter
	dbChapter, err := h.chapterService.GetChapterByNumber(dbNovel.ID, int64(chapterNum))
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	go h.chapterService.UpdateLastReadChapter(userID, dbNovel.ID, dbChapter.ID)

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

	novelStatus := "Completed"
	if dbNovel.IsCompleted == 0 {
		novelStatus = "Ongoing"
	}

	var isRedirect bool = false
	isRedirectStr := pkg.GetAndClearCookie(w, r, comments.CommentRedirectCookie)
	if isRedirectStr != "" {
		isRedirect = true
	}
	chapter := MapToChapterDto(&dbNovel, &dbChapter, prevChapterIntPointer, nextChapterIntPointer)
	metaData := BuildChapterMeta(&dbNovel, chapterNum, novelStatus)

	handlers.GenericHandler(w, r, metaData, chapters.ChapterReader(&chapter, isRedirect))
}

func (h *ChapterHandler) GetChaptersDropDown(w http.ResponseWriter, r *http.Request) {
	novelSlug := chi.URLParam(r, "novelSlug")
	chapterNumberStr := r.URL.Query().Get("chapterNumber")

	chapterNumber, err := strconv.Atoi(chapterNumberStr)
	if err != nil || chapterNumber < 1 {
		handlers.ServerErrorHandler(w, r)
		return
	}

	dbNovel, err := h.chapterService.GetNovelBySlug(novelSlug)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	dbChapters, err := h.chapterService.ListChaptersByNovel(dbNovel.ID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	chaptersList := DbChaptersToChaptersMapper(dbChapters)

	cmp := dropdown.ChapterDropdown(novelSlug, chaptersList, chapterNumber)
	templ.Handler(cmp).ServeHTTP(w, r)
}
