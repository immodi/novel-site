package feedback

import (
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/index"
	"immodi/novel-site/internal/config"
	"immodi/novel-site/internal/http/templates/feedback"
	"net/http"
)

func (h *FeedbackHandler) Feedback(w http.ResponseWriter, r *http.Request) {
	handlers.GenericHandler(w, r, index.BuildHomeMeta(), feedback.GoogleFormSection(config.GoogleFormURL))
}
