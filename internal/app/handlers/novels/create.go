package novels

import (
	"fmt"
	"immodi/novel-site/pkg"
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
		http.Error(w, fmt.Sprintf("Failed to create novel: %v", err), http.StatusInternalServerError)
		return
	}

	// Add default genres
	defaultGenres := pkg.RandomGenres(3)
	for _, genre := range defaultGenres {
		if err := h.novelService.AddGenreToNovel(dbNovel.ID, genre); err != nil {
			http.Error(w, fmt.Sprintf("Failed to add genre %s: %v", genre, err), http.StatusInternalServerError)
			return
		}
	}

	// Add default tags
	defaultTags := pkg.RandomGenres(5)
	for _, tag := range defaultTags {
		if err := h.novelService.AddTagToNovel(dbNovel.ID, tag); err != nil {
			http.Error(w, fmt.Sprintf("Failed to add tag %s: %v", tag, err), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
