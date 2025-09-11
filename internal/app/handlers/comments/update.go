package comments

import (
	"fmt"
	"immodi/novel-site/internal/app/handlers"
	"log"
	"net/http"
	"strconv"
)

func (h *CommentHandler) EditComment(w http.ResponseWriter, r *http.Request) {
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

	novelIDStr := r.FormValue("novelID")
	commentIDStr := r.FormValue("commentID")
	content := r.FormValue("content")

	novelID, err := strconv.Atoi(novelIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to parse novelID")
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

	SetCommentRedirect(w, novelID)
	http.Redirect(w, r, fmt.Sprintf("%s#comments", ref), http.StatusSeeOther)
}
