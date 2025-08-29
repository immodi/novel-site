package privacyservice

import (
	"immodi/novel-site/internal/app"
	homeservice "immodi/novel-site/internal/app/home_service"
	"immodi/novel-site/internal/http/templates/privacy"
	"net/http"
)

func PrivacyHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, homeservice.IndexMetaData, privacy.Privacy())
}
