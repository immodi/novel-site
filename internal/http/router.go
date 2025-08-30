package http

import (
	handlers "immodi/novel-site/internal/app/handlers"
	pathhandlers "immodi/novel-site/internal/app/handlers/paths"
	"immodi/novel-site/internal/app/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	r        *chi.Mux
	handlers *handlers.Handlers
	services *services.Services
}

func (router *Router) NewRouter() *chi.Mux {
	router.r = chi.NewRouter()

	log.Println("Application started at http://localhost:3000")
	router.r.Use(middleware.Logger)

	router.RegisterServices()
	router.RegisterHandlers()
	router.RegisterRoutes()

	return router.r
}

func (router *Router) RegisterRoutes() {
	router.r.Get("/", router.handlers.Home.Index)

	router.r.Get("/novel/{novelName}", router.handlers.Novel.GetNovel)
	router.r.Get("/privacy", pathhandlers.PrivacyHandler)
	router.r.Get("/terms", pathhandlers.TermsHandler)
	router.r.Get("/login", pathhandlers.LoginHandler)
	router.r.Get("/register", pathhandlers.RegisterHandler)
	router.r.Get("/novels", router.redirectToHome())

	// testing routes, should be disabled in production
	router.r.Get("/create-novel", router.handlers.Novel.CreateNovel)
	router.r.Get("/create-chapter/{novelId}", router.handlers.Chapter.CreateChapterWithDefaults)

	router.r.NotFound(router.redirectToHome())
}

func (router *Router) RegisterServices() {
	router.services = services.RegisterServices()
}

func (router *Router) RegisterHandlers() {
	router.handlers = handlers.RegisterHandlers(router.services)
}

func (router *Router) redirectToHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
