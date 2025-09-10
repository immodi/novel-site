package profile

import (
	"immodi/novel-site/internal/app/handlers/novels"
	"log"
	"net/http"
)

func (h *ProfileHandler) RemoveBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		novels.SetErrorMessage(w, "You must be logged in to remove a bookmark.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		novels.SetErrorMessage(w, "Failed to read form data.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("parse form error:", err)
		return
	}

	novelSlug := r.FormValue("slug")
	if novelSlug == "" {
		novels.SetErrorMessage(w, "Invalid novel selected.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dbNovel, err := h.profileService.GetNovelBySlug(novelSlug)
	if err != nil {
		novels.SetErrorMessage(w, "Novel not found.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println(err.Error())
		return
	}

	err = h.profileService.RemoveUserBookmark(userID, dbNovel.ID)
	if err != nil {
		novels.SetErrorMessage(w, "Failed to remove bookmark.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println(err.Error())
		return
	}

	novels.SetSuccessMessage(w, "Bookmark removed successfully!")

	target := r.Referer()
	if target == "" {
		target = "/"
	}

	http.Redirect(w, r, target, http.StatusSeeOther)
}
