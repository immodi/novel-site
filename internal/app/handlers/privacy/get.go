package privacy

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/http/templates/privacy"
	"net/http"
)

func (h *PrivacyHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	handlers.GenericServiceHandler(w, r, index.GetIndexMetaData(), privacy.Privacy())
}
