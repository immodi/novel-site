package app

import (
	"immodi/novel-site/internal/http/templates"
	"immodi/novel-site/internal/http/templates/components"
	"net/http"

	"github.com/a-h/templ"
)

func GenericServiceHandler(
	w http.ResponseWriter,
	r *http.Request,
	data *components.MetaDataStruct,
	cmp templ.Component,
) {
	Render(data, cmp).ServeHTTP(w, r)
}

func Render(data *components.MetaDataStruct, cmp templ.Component) http.HandlerFunc {
	layout := templates.Layout(data, cmp)
	return templ.Handler(layout).ServeHTTP
}
