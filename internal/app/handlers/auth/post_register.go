package auth

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/config"
	"immodi/novel-site/internal/db/repositories"
	sql "immodi/novel-site/internal/db/schema"
	authdtostructs "immodi/novel-site/internal/http/structs/auth"
	"immodi/novel-site/internal/http/templates/auth"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
)

func (h *AuthHandler) PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	terms := r.FormValue("terms") // string: "on" if checked, empty if unchecked

	errors := authenticateRegisterRequest(username, email, password, confirmPassword, terms)

	if len(errors) > 0 {
		dto := authdtostructs.RegisterDTO{
			Errors:   errors,
			Username: username,
			Email:    email,
		}
		handlers.GenericHandler(w, r, BuildAuthMeta("Register"), auth.Register(dto))
		return
	}

	// hash password
	passwordHash, err := pkg.HashPassword(password)
	if err != nil {
		dto := authdtostructs.RegisterDTO{
			Errors:   []string{"Could not store user"},
			Username: username,
			Email:    email,
		}
		handlers.GenericHandler(w, r, BuildAuthMeta("Register"), auth.Register(dto))
		return
	}

	// create the user
	userId, err := h.authService.RegisterUser(repositories.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         string(sql.UserRoleUser),
		CreatedAt:    pkg.GetCurrentTimeRFC3339(),
	})

	if err != nil {
		log.Println(err.Error())
		dto := authdtostructs.RegisterDTO{
			Errors:   []string{err.Error()},
			Username: username,
			Email:    email,
		}
		handlers.GenericHandler(w, r, BuildAuthMeta("Register"), auth.Register(dto))
		return
	}

	token, err := GenerateToken(userId, string(sql.UserRoleUser), DefaultJwtDuration)
	if err != nil {
		dto := authdtostructs.RegisterDTO{
			Errors:   []string{"Could not generate token"},
			Username: username,
			Email:    email,
		}
		handlers.GenericHandler(w, r, BuildAuthMeta("Register"), auth.Register(dto))
		return
	}

	setAuthCookie(w, token, DefaultJwtDuration, config.IsProduction)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
