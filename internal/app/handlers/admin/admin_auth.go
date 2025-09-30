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
			Error: "method not allowed",
		})
		return
	}

	var req admin.AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Error: "invalid request body",
		})
		return
	}
	defer r.Body.Close()

	error := authenticateLoginRequest(req.Email, req.Password)

	if len(error) > 0 {
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Error: "invalid request body",
		})
		return
	}

	user, err := h.authService.LoginUserWithEmail(req.Email, req.Password)
	if err != nil {
		error := "Invalid email or password"
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Error: error,
		})
		return
	}

	if user.Role != string(sql.UserRoleAdmin) {
		error := "Your'e not an admin buddy"
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Error: error,
		})
		return
	}

	token, err := pkg.GenerateToken(user.ID, user.Role, pkg.DefaultJwtDuration)
	if err != nil {
		error := "Could not generate token"
		handlers.WriteJSON(w, http.StatusBadRequest, admin.AdminLoginResponse{
			Error: error,
		})
		return
	}

	handlers.WriteJSON(w, http.StatusOK, admin.AdminLoginResponse{
		Token:      token,
		Username:   user.Username,
		CoverImage: user.Image,
		Error:      "",
	})
}
