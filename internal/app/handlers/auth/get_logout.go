package auth

import (
	"immodi/novel-site/internal/config"
	"net/http"
)

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	removeAuthCookie(w, config.IsProduction)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
