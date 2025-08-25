package http

import (
	"immodi/novel-site/internal/http/routes"
	"net/http"

	"github.com/a-h/templ"
)

func Render(title string, cmp templ.Component) http.HandlerFunc {
	layout := routes.Layout(title, cmp)
	return templ.Handler(layout).ServeHTTP
}
