package admin

import (
	"net/http"
	"strconv"
	"strings"

	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/internal/http/payloads/admin"
	"immodi/novel-site/pkg"

	"github.com/go-chi/chi/v5"
)

func (h *AdminHandler) AdminGetAllNovelChapters(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelChaptersResponse{
			Chapters: nil,
			Error:    "missing Authorization header",
		})
		return
	}

	novelIDStr := chi.URLParam(r, "novelID")
	novelID, err := strconv.Atoi(novelIDStr)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelChaptersResponse{
			Chapters: nil,
			Error:    "invalid novelID",
		})
		return
	}

	// expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelChaptersResponse{
			Chapters: nil,
			Error:    "invalid Authorization header format",
		})
		return
	}
	token := parts[1]

	userID, err := pkg.GetUserIDFromToken(token)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelChaptersResponse{
			Chapters: nil,
			Error:    "invalid token",
		})
		return
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelChaptersResponse{
			Chapters: nil,
			Error:    "could not get the user from the token",
		})
		return
	}

	if user.Role != string(sql.UserRoleAdmin) {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelChaptersResponse{
			Chapters: nil,
			Error:    "user is not an admin",
		})
		return
	}

	chapters, err := h.chapterService.ListChaptersByNovel(int64(novelID))
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelChaptersResponse{
			Chapters: nil,
			Error:    "could not get all chapters",
		})
		return
	}

	adminPanelChapters := DbChaptersToAdminPanelChaptersMapper(chapters)
	handlers.WriteJSON(w, http.StatusOK, admin.AdminGetAllNovelChaptersResponse{
		Chapters: adminPanelChapters,
		Error:    "",
	})
}
