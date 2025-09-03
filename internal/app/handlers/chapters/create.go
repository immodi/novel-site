package chapters

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *ChapterHandler) CreateChapterWithDefaults(w http.ResponseWriter, r *http.Request) {
	novelIDString := chi.URLParam(r, "novelId")
	novelID, err := strconv.ParseInt(novelIDString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid novel ID", http.StatusBadRequest)
		return
	}

	_, err = h.chapterService.CreateChapterWithDefaults(novelID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create chapter: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/novel", http.StatusSeeOther)
}
