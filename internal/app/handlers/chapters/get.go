package chapters

import (
	"immodi/novel-site/internal/app/handlers"
	dropdown "immodi/novel-site/internal/http/components/chapters"
	"immodi/novel-site/internal/http/templates/chapters"
	"immodi/novel-site/pkg"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func (h *ChapterHandler) ReadChapter(w http.ResponseWriter, r *http.Request) {
	novelName := pkg.SlugToTitle(chi.URLParam(r, "novelName"))
	chapterNumStr := pkg.SlugToTitle(chi.URLParam(r, "chapterNumber"))

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil || chapterNum < 1 {
		handlers.ServerErrorHandler(w, r)
		return
	}

	dbNovel, err := h.chapterService.GetNovelByNameLike(novelName)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

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

	chapter := BuildChapterPage(dbNovel, dbChapter, prevChapterIntPointer, nextChapterIntPointer)
	metaData := BuildChapterMeta(dbNovel, chapterNum, novelStatus)

	handlers.GenericServiceHandler(w, r, metaData, chapters.ChapterReader(chapter))
}

func (h *ChapterHandler) GetChaptersDropDown(w http.ResponseWriter, r *http.Request) {
	novelName := pkg.SlugToTitle(chi.URLParam(r, "novelName"))

	dbNovel, err := h.chapterService.GetNovelByNameLike(novelName)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	dbChapters, err := h.chapterService.ListChaptersByNovel(dbNovel.ID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		return
	}

	chaptersList := dbChaptersToChapters(dbChapters)

	cmp := dropdown.ChapterDropdown(pkg.TitleToSlug(novelName), chaptersList)
	templ.Handler(cmp).ServeHTTP(w, r)

}
