package chaptercomments

import (
	"fmt"
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/comments"
	"log"
	"net/http"
	"strconv"
)

func (h *ChapterCommentHandler) EditComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	chapterIDStr := r.FormValue("chapterID")
	commentIDStr := r.FormValue("commentID")
	content := r.FormValue("content")

	chapterID, err := strconv.Atoi(chapterIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to parse chapterID")
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	comment, err := h.commentService.GetCommentById(commentID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	if comment.UserID != userID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = h.commentService.UpdateComment(commentID, content)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	ref := r.Referer()
	if ref == "" {
		ref = "/"
	}

	comments.SetCommentRedirect(w, chapterID)
	http.Redirect(w, r, fmt.Sprintf("%s#comments", ref), http.StatusSeeOther)
}
