package novels

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/pkg"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *NovelHandler) CreateNovelWithDefaults(w http.ResponseWriter, r *http.Request) {
	novelName := chi.URLParam(r, "novelName")
	novelStatusStr := chi.URLParam(r, "novelStatus")
	novelStatus := false
	if novelStatusStr == "1" {
		novelStatus = true
	}

	dbNovel, err := h.novelService.CreateNovelWithDefaults(novelName, novelStatus)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	// Add default genres
	defaultGenres := pkg.RandomGenres(3)
	for _, genre := range defaultGenres {
		if err := h.novelService.AddGenreToNovel(dbNovel.ID, genre); err != nil {
			handlers.ServerErrorHandler(w, r)
			log.Println(err.Error())
			return
		}
	}

	// Add default tags
	defaultTags := pkg.RandomGenres(5)
	for _, tag := range defaultTags {
		if err := h.novelService.AddTagToNovel(dbNovel.ID, tag); err != nil {
			handlers.ServerErrorHandler(w, r)
			log.Println(err.Error())
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
