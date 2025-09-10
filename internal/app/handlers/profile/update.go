package profile

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/db/repositories"
	"log"
	"net/http"
)

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println("failed to parse form:", err)
		return
	}

	// username := r.FormValue("username")
	image := r.FormValue("image")

	_, err := h.profileService.UpdateUserPartial(repositories.UpdateUserPartialParams{
		ID: userID,
		// Username:     username,
		Image:        image,
		PasswordHash: nil,
		Role:         nil,
		CreatedAt:    nil,
	})

	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println("failed to update user:", err)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
