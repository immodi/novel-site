package admin

import (
	"encoding/json"
	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/internal/http/payloads/admin"
	"immodi/novel-site/pkg"
	"net/http"
)

func (h *AdminHandler) AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.WriteJSON(w, http.StatusMethodNotAllowed, admin.AdminLoginResponse{
			Token:  "",
			Errors: []string{"method not allowed"},
		})
		return
	}

	var req admin.AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: []string{"invalid request body"},
		})
		return
	}
	defer r.Body.Close()

	errors := authenticateLoginRequest(req.Email, req.Password)

	if len(errors) > 0 {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: errors,
		})
		return
	}

	user, err := h.authService.LoginUserWithEmail(req.Email, req.Password)
	if err != nil {
		errors := []string{"Invalid email or password"}
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: errors,
		})
		return
	}

	if user.Role != string(sql.UserRoleAdmin) {
		errors := []string{"Your'e not an admin buddy"}
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: errors,
		})
		return
	}

	token, err := pkg.GenerateToken(user.ID, user.Role, pkg.DefaultJwtDuration)
	if err != nil {
		errors := []string{"Could not generate token"}
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Token:  "",
			Errors: errors,
		})
		return
	}

	handlers.WriteJSON(w, http.StatusOK, admin.AdminLoginResponse{
		Token:  token,
		Errors: nil,
	})
}
