package handlers

import (
	"immodi/novel-site/internal/http/templates/terms"
	"net/http"
)

func TermsHandler(w http.ResponseWriter, r *http.Request) {
	GenericServiceHandler(w, r, IndexMetaData, terms.Terms())
}
