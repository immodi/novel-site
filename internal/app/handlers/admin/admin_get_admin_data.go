package admin

import (
	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/internal/http/payloads/admin"
	"immodi/novel-site/pkg"
	"net/http"
	"strings"
)

func (h *AdminHandler) AdminGetAdminData(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAdminDataResponse{
			Error: "missing Authorization header",
		})
		return
	}

	// expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAdminDataResponse{
			Error: "invalid Authorization header format",
		})
		return
	}
	token := parts[1]

	userID, err := pkg.GetUserIDFromToken(token)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAdminDataResponse{
			Error: "invalid token",
		})
		return
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAdminDataResponse{
			Error: "could not get the user from the token",
		})
		return
	}

	if user.Role != string(sql.UserRoleAdmin) {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminGetAdminDataResponse{
			Error: "user is not an admin",
		})
		return
	}

	handlers.WriteJSON(w, http.StatusOK, admin.AdminGetAdminDataResponse{
		Username:   user.Username,
		CoverImage: user.Image,
	})
}
