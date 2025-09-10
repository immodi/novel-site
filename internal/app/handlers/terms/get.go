package terms

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/http/templates/terms"
	"net/http"
)

func (h *TermsHandler) Terms(w http.ResponseWriter, r *http.Request) {
	handlers.GenericHandler(w, r, index.BuildHomeMeta(), terms.Terms())
}
