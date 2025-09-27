package admin

import (
	"encoding/json"
	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/internal/http/payloads/admin"
	"immodi/novel-site/pkg"
	"net/http"
)

func (h *AdminHandler) AdminGetAllUsers(w http.ResponseWriter, r *http.Request) {
	var req admin.AdminGetAllUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: []string{"invalid request body"},
		})
		return
	}
	defer r.Body.Close()

	userID, err := pkg.GetUserIDFromToken(req.Token)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: []string{"invalid token"},
		})
		return
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: []string{"coudlnt get the user from the admin token"},
		})
		return
	}

	if user.Role != string(sql.UserRoleAdmin) {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: []string{"user is not an admin"},
		})
		return
	}

	users, err := h.authService.GetAllUsers()
	if err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: []string{"could not get all users"},
		})
		return
	}

	adminPanelUsers := DbUsersToAdminPanelUsersMapper(users)
	handlers.WriteJSON(w, http.StatusOK, admin.AdminGetAllUsersResponse{
		Users: adminPanelUsers,
	})
}
