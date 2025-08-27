package termsservice

import (
	"immodi/novel-site/internal/app"
	"immodi/novel-site/internal/http/routes/terms"
	"net/http"
)

func TermsHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, "Terms", terms.Terms())
}
