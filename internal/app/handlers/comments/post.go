package comments

import (
	"fmt"
	"immodi/novel-site/internal/app/handlers"
	"log"
	"net/http"
	"strconv"
)

func (h *CommentHandler) PostComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	novelIDStr := r.FormValue("novelID")
	parentCommentIDStr := r.FormValue("parentCommentID")
	content := r.FormValue("content")

	novelID, err := strconv.Atoi(novelIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error())
		return
	}

	var parentCommentIDPointer *int
	parentCommentID, err := strconv.Atoi(parentCommentIDStr)
	if err == nil {
		parentCommentIDPointer = &parentCommentID
	} else {
		parentCommentIDPointer = nil
	}

	_, err = h.commentService.CreateComment(novelID, int(user.ID), parentCommentIDPointer, content)
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
