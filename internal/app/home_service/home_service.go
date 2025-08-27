package homeservice

import (
	"immodi/novel-site/internal/app"
	"immodi/novel-site/internal/http/routes/index"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.GenericServiceHandler(w, r, "Home", index.Index())
}
