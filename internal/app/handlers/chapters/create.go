package chapters

import (
	"immodi/novel-site/internal/app/handlers"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *ChapterHandler) CreateChapterWithDefaults(w http.ResponseWriter, r *http.Request) {
	novelIDString := chi.URLParam(r, "novelId")
	novelID, err := strconv.ParseInt(novelIDString, 10, 64)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	_, err = h.chapterService.CreateChapterWithDefaults(novelID)

	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
