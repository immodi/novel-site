package http

import (
	authservice "immodi/novel-site/internal/app/auth_service"
	homeservice "immodi/novel-site/internal/app/home_service"
	novelservice "immodi/novel-site/internal/app/novel_service"
	privacyservice "immodi/novel-site/internal/app/privacy_service"
	termsservice "immodi/novel-site/internal/app/terms_service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	r *chi.Mux
}

func (router *Router) NewRouter() *chi.Mux {
	router.r = chi.NewRouter()

	log.Println("Application started at http://localhost:3000")
	router.r.Use(middleware.Logger)
	router.RegisterRoutes()

	return router.r
}

func (router *Router) RegisterRoutes() {
	router.r.Get("/", homeservice.HomeHandler)

	router.r.Get("/novel", novelservice.NovelHandler)
	router.r.Get("/privacy", privacyservice.PrivacyHandler)
	router.r.Get("/terms", termsservice.TermsHandler)
	router.r.Get("/login", authservice.LoginHandler)
	router.r.Get("/register", authservice.RegisterHandler)
	router.r.Get("/novels", router.redirectToHome())

	router.r.NotFound(router.redirectToHome())
}

func (router *Router) redirectToHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
