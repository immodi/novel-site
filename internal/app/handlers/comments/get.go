package comments

import (
	"immodi/novel-site/internal/app/handlers"
	comment_component "immodi/novel-site/internal/http/components/novels"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

func (h *CommentHandler) Comments(w http.ResponseWriter, r *http.Request) {
	userID := IsAuthedUser(r)
	novelIDStr := r.URL.Query().Get("novelID")

	novelID, err := strconv.Atoi(novelIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	dbComments, err := h.commentService.GetCommentsByNovel(novelID)

	comments := DbCommentsToCommentDtoMapper(dbComments, h.commentService)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}
	cmp := comment_component.CommentsComponent(comments, int(userID), novelID)
	templ.Handler(cmp).ServeHTTP(w, r)
}
