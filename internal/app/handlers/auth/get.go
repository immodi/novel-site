package auth

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/http/templates/auth"
	"net/http"
)

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	handlers.GenericServiceHandler(w, r, getAuthMetaData("Login"), auth.Login())

}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	handlers.GenericServiceHandler(w, r, getAuthMetaData("Register"), auth.Register())
}
