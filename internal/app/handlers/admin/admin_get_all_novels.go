package admin

import (
	"net/http"
	"strings"

	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/internal/http/payloads/admin"
	"immodi/novel-site/pkg"
)

func (h *AdminHandler) AdminGetAllNovels(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelsResponse{
			Novels: nil,
			Error:  "missing Authorization header",
		})
		return
	}

	// expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelsResponse{
			Novels: nil,
			Error:  "invalid Authorization header format",
		})
		return
	}
	token := parts[1]

	userID, err := pkg.GetUserIDFromToken(token)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelsResponse{
			Novels: nil,
			Error:  "invalid token",
		})
		return
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelsResponse{
			Novels: nil,
			Error:  "could not get the user from the token",
		})
		return
	}

	if user.Role != string(sql.UserRoleAdmin) {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelsResponse{
			Novels: nil,
			Error:  "user is not an admin",
		})
		return
	}

	novels, err := h.indexService.ListAllNovels()
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAllNovelsResponse{
			Novels: nil,
			Error:  "could not get all novels",
		})
		return
	}

	adminPanelNovels := DbNovelsToAdminPanelNovelsMapper(novels)
	handlers.WriteJSON(w, http.StatusOK, admin.AdminGetAllNovelsResponse{
		Novels: adminPanelNovels,
		Error:  "",
	})
}
