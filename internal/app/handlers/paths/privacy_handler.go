package handlers

import (
	"immodi/novel-site/internal/http/templates/privacy"
	"net/http"
)

func PrivacyHandler(w http.ResponseWriter, r *http.Request) {
	GenericServiceHandler(w, r, IndexMetaData, privacy.Privacy())
}
