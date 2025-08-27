package app

import (
	"immodi/novel-site/internal/http/routes"
	"net/http"

	"github.com/a-h/templ"
)

func GenericServiceHandler(w http.ResponseWriter, r *http.Request, title string, cmp templ.Component) {
	Render(title, cmp).ServeHTTP(w, r)
}

func Render(title string, cmp templ.Component) http.HandlerFunc {
	layout := routes.Layout(title, cmp)
	return templ.Handler(layout).ServeHTTP
}
