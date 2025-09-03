package http

import (
	"immodi/novel-site/internal/http/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	r        *chi.Mux
	handlers *utils.Handlers
	services *utils.Services
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
	router.r.Get("/novel/{novelName}/chapters", router.handlers.Chapter.GetChaptersDropDown)
	router.r.Get("/novel/{novelName}/chapter-{chapterNumber}", router.handlers.Chapter.ReadChapter)

	router.r.Get("/privacy", router.handlers.Privacy.Privacy)
	router.r.Get("/terms", router.handlers.Terms.Terms)
	router.r.Get("/login", router.handlers.Auth.LoginHandler)
	router.r.Get("/register", router.handlers.Auth.RegisterHandler)
	router.r.Get("/novels", router.redirectToHome())

	// testing routes, should be disabled in production
	router.r.Get("/create-novel/{novelName}/{novelStatus}", router.handlers.Novel.CreateNovelWithDefaults)
	router.r.Get("/create-chapter/{novelId}", router.handlers.Chapter.CreateChapterWithDefaults)

	router.r.NotFound(router.redirectToHome())
}

func (router *Router) RegisterServices() {
	router.services = utils.RegisterServices()
}

func (router *Router) RegisterHandlers() {
	router.handlers = utils.RegisterHandlers(router.services)
}

func (router *Router) redirectToHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
