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
	statusCode ...int,
) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	authHeaderEntry := "Login"
	_, err := r.Cookie("auth_token")
	if err == nil {
		authHeaderEntry = "Profile"
	}

	headers := []string{"Filter", authHeaderEntry}

	Render(data, &indexdtostructs.LayoutData{Headers: headers}, cmp, code).ServeHTTP(w, r)
}

func Render(metaData *indexdtostructs.MetaDataStruct, data *indexdtostructs.LayoutData, cmp templ.Component, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		layout := templates.Layout(metaData, data, cmp)

		w.WriteHeader(code)
		templ.Handler(layout).ServeHTTP(w, r)
	}
}
