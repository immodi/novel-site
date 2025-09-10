package auth

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/config"
	authdtostructs "immodi/novel-site/internal/http/structs/auth"
	"immodi/novel-site/internal/http/templates/auth"
	"net/http"
)

func (h *AuthHandler) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	errors := authenticateLoginRequest(email, password)

	if len(errors) > 0 {
		dto := authdtostructs.LoginDTO{
			Errors: errors,
			Email:  email,
		}
		handlers.GenericHandler(w, r, BuildAuthMeta("Login"), auth.Login(dto))
		return
	}

	user, err := h.authService.LoginUserWithEmail(email, password)
	if err != nil {
		dto := authdtostructs.LoginDTO{Errors: []string{"Invalid email or password"}, Email: email}
		handlers.GenericHandler(w, r, BuildAuthMeta("Login"), auth.Login(dto))
	}

	token, err := GenerateToken(user.ID, user.Role, DefaultJwtDuration)
	if err != nil {
		dto := authdtostructs.LoginDTO{Errors: []string{"Could not generate token"}, Email: email}
		handlers.GenericHandler(w, r, BuildAuthMeta("Login"), auth.Login(dto))
		return
	}

	setAuthCookie(w, token, DefaultJwtDuration, config.IsProduction)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
