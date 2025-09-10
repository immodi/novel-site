package handlers

import (
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	"immodi/novel-site/internal/http/templates"
	"net/http"

	"github.com/a-h/templ"
)

func GenericHandler(
	w http.ResponseWriter,
	r *http.Request,
	data *indexdtostructs.MetaDataStruct,
	cmp templ.Component,
) {
	authHeaderEntry := "Login"
	_, err := r.Cookie("auth_token")
	if err == nil {
		authHeaderEntry = "Profile"
	}

	headers := []string{"Novels", authHeaderEntry}

	Render(data, &indexdtostructs.LayoutData{Headers: headers}, cmp).ServeHTTP(w, r)
}

func Render(metaData *indexdtostructs.MetaDataStruct, data *indexdtostructs.LayoutData, cmp templ.Component) http.HandlerFunc {
	layout := templates.Layout(metaData, data, cmp)
	return templ.Handler(layout).ServeHTTP
}
