package http

import (
	"fmt"
	"immodi/novel-site/internal/http/routes/about"
	"immodi/novel-site/internal/http/routes/index"
	"immodi/novel-site/internal/http/routes/novels"
	"immodi/novel-site/internal/http/routes/privacy"
	"immodi/novel-site/internal/http/routes/terms"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	r *chi.Mux
}

func (router *Router) NewRouter() *chi.Mux {
	router.r = chi.NewRouter()

	router.r.Use(middleware.Logger)
	router.RegisterRoutes()

	return router.r
}

func (router *Router) RegisterRoutes() {
	router.r.Get("/", Render("Home", index.Index()))

	chapters := []novels.Chapter{}
	for i := 1; i <= 100; i++ {
		chapters = append(chapters, novels.Chapter{
			Title:  fmt.Sprintf("Chapter %d", i),
			Number: i,
		})
	}

	router.r.Get("/novel", Render("Test Novel", novels.NovelInfo(novels.Novel{
		Name:        "Test Novel",
		Description: "Test Novel Description",
		Author:      "Test Author",
		Genres:      []string{"Test Genre", "Test Genre 2"},
		Status:      "Ongoing",
		CoverImage:  "https://dummyimage.com/600x400/000/fff",
		Chapters:    chapters,
	})))

	router.r.Get("/privacy", Render("Privacy", privacy.Privacy()))
	router.r.Get("/terms", Render("Terms of Service", terms.Terms()))
	router.r.Get("/about", Render("About", about.About()))

	router.r.Get("/novels", router.redirectToHome())

	router.r.NotFound(router.redirectToHome())
}

func (router *Router) redirectToHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
