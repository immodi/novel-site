package privacyservice

import (
	"immodi/novel-site/internal/app"
	"immodi/novel-site/internal/http/routes/privacy"
	"net/http"
)

func PrivacyHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, "Privacy", privacy.Privacy())
}
