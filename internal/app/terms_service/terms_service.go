package termsservice

import (
	"immodi/novel-site/internal/app"
	homeservice "immodi/novel-site/internal/app/home_service"
	"immodi/novel-site/internal/http/templates/terms"
	"net/http"
)

func TermsHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, homeservice.IndexMetaData, terms.Terms())
}
