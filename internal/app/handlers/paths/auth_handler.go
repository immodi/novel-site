package handlers

import (
	"fmt"
	"immodi/novel-site/internal/http/templates/auth"
	"immodi/novel-site/internal/http/templates/components"
	"net/http"
)

func getAuthMetaData(title string) *components.MetaDataStruct {
	return &components.MetaDataStruct{
		Title:       fmt.Sprintf("%s - %s", title, components.SITE_NAME),
		IsRendering: false,
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	GenericServiceHandler(w, r, getAuthMetaData("Login"), auth.Login())
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	GenericServiceHandler(w, r, getAuthMetaData("Register"), auth.Register())
}
