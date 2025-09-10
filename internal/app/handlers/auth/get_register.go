package auth

import (
	"immodi/novel-site/internal/app/handlers"
	authdtostructs "immodi/novel-site/internal/http/structs/auth"
	"immodi/novel-site/internal/http/templates/auth"
	"net/http"
)

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	handlers.GenericHandler(w, r, BuildAuthMeta("Register"), auth.Register(authdtostructs.RegisterDTO{}))
}
