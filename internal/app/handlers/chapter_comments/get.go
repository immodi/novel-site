package chaptercomments

import (
	"immodi/novel-site/internal/app/handlers"
	comment_component "immodi/novel-site/internal/http/components/chapters"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

func (h *ChapterCommentHandler) Comments(w http.ResponseWriter, r *http.Request) {
	userID := pkg.IsAuthedUser(r)
	chapterIDStr := r.URL.Query().Get("chapterID")

	chapterID, err := strconv.Atoi(chapterIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	dbComments, err := h.commentService.GetCommentsByChapter(chapterID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	commentsIDs := make([]int64, len(dbComments))
	for i, c := range dbComments {
		commentsIDs[i] = c.ID
	}

	dbCommentsReactions, err := h.commentService.GetUserReactionForComments(int(userID), commentsIDs)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}
	comments := DbCommentsToCommentDtoMapper(dbComments, dbCommentsReactions, int(userID), h.commentService)

	// dbCommentsWithUserReactions, err := h.commentService.GetCommentsByChapterWithUserReactions(int(userID), chapterID)
	// if err != nil {
	// 	handlers.ServerErrorHandler(w, r)
	// 	log.Println(err.Error())
	// 	return
	// }

	// comments := DbCommentsWithReactionsToCommentDtoMapper(dbCommentsWithUserReactions, h.commentService)

	cmp := comment_component.ChapterCommentsComponent(comments, int(userID), chapterID)
	templ.Handler(cmp).ServeHTTP(w, r)
}
