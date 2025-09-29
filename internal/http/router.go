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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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
	router.r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, time.Minute))

		r.Get("/", router.handlers.Home.Index)

		r.Get("/novel/{novelSlug}", router.handlers.Novel.GetNovel)
		r.Get("/novel/{novelSlug}/chapters", router.handlers.Chapter.GetChaptersDropDown)
		r.Get("/novel/{novelSlug}/chapter-{chapterNumber}", router.handlers.Chapter.ReadChapter)

		r.Get("/comments", router.handlers.Comment.Comments)
		r.Get("/chapter-comments", router.handlers.ChapterComment.Comments)

		r.Get("/search/{novelName}", router.handlers.Search.SearchNovel)
		r.Get("/sort/{collection}", router.handlers.Search.SortNovelsByCollection)
		r.Get("/genre/{genre}", router.handlers.Search.SortNovelsByGenres)
		r.Get("/tag/{tag}", router.handlers.Search.SortNovelsByTags)
		r.Get("/author/{author}", router.handlers.Search.SortNovelsByAuthor)
		r.Get("/filter-results", router.handlers.Search.SortByFilter)

		r.Get("/filter", router.handlers.Filter.Filter)
		r.Get("/filter/tags", router.handlers.Filter.FilterTags)
		r.Get("/filter/tags-excluded", router.handlers.Filter.FilterExcludedTags)

		r.Get("/privacy", router.handlers.Privacy.Privacy)
		r.Get("/terms", router.handlers.Terms.Terms)

		r.Get("/auth/google/login", router.handlers.Auth.GoogleAuth)
		r.Get("/auth/google/callback", router.handlers.Auth.GoogleCallback)

		r.Get("/login", router.handlers.Auth.LoginHandler)
		r.Post("/login", router.handlers.Auth.PostLoginHandler)
		r.Get("/logout", router.handlers.Auth.LogoutHandler)

		r.Get("/register", router.handlers.Auth.RegisterHandler)
		r.Post("/register", router.handlers.Auth.PostRegisterHandler)

		if !config.IsProduction {
			r.Get("/create-novel/{novelName}/{novelStatus}", router.handlers.Novel.CreateNovelWithDefaults)
			r.Get("/create-chapter/{novelId}", router.handlers.Chapter.CreateChapterWithDefaults)
		}

		r.Group(func(r chi.Router) {
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

		r.Post("/admin/login", router.handlers.Admin.AdminLoginHandler)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.ApiAdminRoleMiddleware())
			r.Get("/admin/users", router.handlers.Admin.AdminGetAllUsers)
			r.Get("/admin/novels", router.handlers.Admin.AdminGetAllNovels)
		})

		r.Handle("/static/*", router.serveStatic("static"))
		r.Get("/robots.txt", router.serveStaticAsset("robots.txt"))
		r.Get("/favicon.ico", router.serveStaticAsset("logo/favicon.ico"))

		r.Get("/sitemap.xml", router.handlers.Sitemap.MainSiteMap)
		r.Get("/sitemaps/home.xml", router.handlers.Sitemap.HomeSiteMap)
		r.Get("/sitemaps/novels.xml", router.handlers.Sitemap.NovelsSiteMap)
		r.Get("/sitemaps/genres.xml", router.handlers.Sitemap.GenresSiteMap)
		r.Get("/sitemaps/tags.xml", router.handlers.Sitemap.TagsSiteMap)

		r.Get("/novels", router.redirectToHome())
		r.NotFound(handlers.NotFoundHandler)

	})

	router.r.Group(func(r chi.Router) {
		r.Use(middlewares.LocalOnlyMiddleware)

		r.Post("/load/novel", router.handlers.Load.LoadNovel)
		r.Post("/load/chapters", router.handlers.Load.LoadChapter)
		r.Post("/load/last-chapter/id", router.handlers.Load.GetLastChapterById)
		r.Post("/load/last-chapter/name", router.handlers.Load.GetLastChapterByName)
		r.Post("/load/append-chapters", router.handlers.Load.AppendChapters)
	})

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
		filePath := filepath.Join("static", assetName)
		info, err := os.Stat(filePath)
		if os.IsNotExist(err) || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week

		http.ServeFile(w, r, filePath)
	}
}
