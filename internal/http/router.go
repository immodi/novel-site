package http

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/middlewares"
	"immodi/novel-site/internal/config"
	"immodi/novel-site/internal/http/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	router.RegisterMiddlewares()
	router.RegisterServices()
	router.RegisterHandlers()
	router.RegisterRoutes()

	return router.r
}

func (router *Router) RegisterRoutes() {
	router.r.Get("/", router.handlers.Home.Index)

	router.r.Get("/novel/{novelSlug}", router.handlers.Novel.GetNovel)
	router.r.Get("/novel/{novelSlug}/chapters", router.handlers.Chapter.GetChaptersDropDown)
	router.r.Get("/novel/{novelSlug}/chapter-{chapterNumber}", router.handlers.Chapter.ReadChapter)

	router.r.Get("/comments", router.handlers.Comment.Comments)
	router.r.Get("/chapter-comments", router.handlers.ChapterComment.Comments)

	router.r.Get("/search/{novelName}", router.handlers.Search.SearchNovel)
	router.r.Get("/sort/{collection}", router.handlers.Search.SortNovelsByCollection)
	router.r.Get("/genre/{genre}", router.handlers.Search.SortNovelsByGenres)
	router.r.Get("/tag/{tag}", router.handlers.Search.SortNovelsByTags)
	router.r.Get("/author/{author}", router.handlers.Search.SortNovelsByAuthor)

	router.r.Get("/privacy", router.handlers.Privacy.Privacy)
	router.r.Get("/terms", router.handlers.Terms.Terms)

	router.r.Get("/login", router.handlers.Auth.LoginHandler)
	router.r.Post("/login", router.handlers.Auth.PostLoginHandler)
	router.r.Get("/logout", router.handlers.Auth.LogoutHandler)

	router.r.Get("/register", router.handlers.Auth.RegisterHandler)
	router.r.Post("/register", router.handlers.Auth.PostRegisterHandler)

	router.r.Post("/load/novel", router.handlers.Load.LoadNovel)
	router.r.Post("/load/chapters", router.handlers.Load.LoadChapter)
	router.r.Post("/load/last-chapter/id", router.handlers.Load.GetLastChapterById)
	router.r.Post("/load/last-chapter/name", router.handlers.Load.GetLastChapterByName)
	router.r.Post("/load/append-chapters", router.handlers.Load.AppendChapters)

	if !config.IsProduction {
		router.r.Get("/create-novel/{novelName}/{novelStatus}", router.handlers.Novel.CreateNovelWithDefaults)
		router.r.Get("/create-chapter/{novelId}", router.handlers.Chapter.CreateChapterWithDefaults)
	}

	router.r.Group(func(r chi.Router) {
		r.Use(middlewares.RoleMiddleware("user"))

		r.Get("/profile", router.handlers.Profile.Profile)
		r.Post("/profile", router.handlers.Profile.UpdateProfile)

		r.Post("/bookmark", router.handlers.Profile.PostBookmark)
		r.Post("/bookmark-remove", router.handlers.Profile.RemoveBookmark)

		r.Post("/comment", router.handlers.Comment.PostComment)
		r.Post("/comment/edit", router.handlers.Comment.EditComment)
		r.Post("/comment/reaction", router.handlers.Comment.PostReact)

		r.Post("/chapter-comments", router.handlers.ChapterComment.PostComment)
		r.Post("/chapter-comments/edit", router.handlers.ChapterComment.EditComment)
		r.Post("/chapter-comments/reaction", router.handlers.ChapterComment.PostReact)

	})

	router.r.Handle("/static/*", router.serveStatic("static"))
	router.r.Get("/robots.txt", router.serveStaticAsset("robots.txt"))

	router.r.Get("/sitemap.xml", router.handlers.Sitemap.MainSiteMap)
	router.r.Get("/sitemaps/novels.xml", router.handlers.Sitemap.NovelsSiteMap)
	router.r.Get("/sitemaps/genres.xml", router.handlers.Sitemap.GenresSiteMap)
	router.r.Get("/sitemaps/tags.xml", router.handlers.Sitemap.TagsSiteMap)

	router.r.Get("/novels", router.redirectToHome())
	router.r.NotFound(handlers.NotFoundHandler)

}

func (router *Router) RegisterServices() {
	router.services = utils.RegisterServices()
}

func (router *Router) RegisterHandlers() {
	router.handlers = utils.RegisterHandlers(router.services)
}

func (router *Router) RegisterMiddlewares() {
	router.r.Use(middleware.Logger)
}

func (router *Router) redirectToHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (router *Router) serveStatic(dir string) http.HandlerFunc {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("failed to resolve static dir: %v", err)
	}

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(absDir)))

	return func(w http.ResponseWriter, r *http.Request) {
		// Prevent directory listing
		filePath := filepath.Join(absDir, r.URL.Path[len("/static/"):])
		if info, err := os.Stat(filePath); os.IsNotExist(err) || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		// Set Cache-Control header: 1 week
		w.Header().Set("Cache-Control", "public, max-age=604800") // 604800 seconds = 7 days

		fs.ServeHTTP(w, r)
	}
}

func (router *Router) serveStaticAsset(assetName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Absolute path to the robots.txt file in your static folder
		filePath := filepath.Join("static", assetName) // adjust if needed
		info, err := os.Stat(filePath)
		if os.IsNotExist(err) || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		// Optional: set cache control
		w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week

		http.ServeFile(w, r, filePath)
	}
}
