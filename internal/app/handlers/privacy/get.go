package privacy

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/http/templates/privacy"
	"net/http"
)

func (h *PrivacyHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	handlers.GenericHandler(w, r, index.BuildHomeMeta(), privacy.Privacy())
}
