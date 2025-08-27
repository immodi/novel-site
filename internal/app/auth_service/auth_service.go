package authservice

import (
	"immodi/novel-site/internal/app"
	"immodi/novel-site/internal/http/routes/auth"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, "Login", auth.Login())
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, "Register", auth.Register())
}
