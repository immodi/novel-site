package handlers

import (
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	"immodi/novel-site/internal/http/templates"
	"net/http"

	"github.com/a-h/templ"
)

func GenericServiceHandler(
	w http.ResponseWriter,
	r *http.Request,
	data *indexdtostructs.MetaDataStruct,
	cmp templ.Component,
) {
	Render(data, cmp).ServeHTTP(w, r)
}

func Render(data *indexdtostructs.MetaDataStruct, cmp templ.Component) http.HandlerFunc {
	layout := templates.Layout(data, cmp)
	return templ.Handler(layout).ServeHTTP
}
