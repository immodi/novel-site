package chaptercomments

import (
	"fmt"
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/app/handlers/comments"
	sql "immodi/novel-site/internal/db/schema"
	"log"
	"net/http"
	"strconv"
)

func (h *ChapterCommentHandler) PostReact(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to get user")
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to parse form")
		return
	}
	defer r.Body.Close()

	chapterIDStr := r.FormValue("chapterID")
	commentIDStr := r.FormValue("commentID")
	reactionStr := r.FormValue("reaction")

	chapterID, err := strconv.Atoi(chapterIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to parse chapterID")
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to parse commentID")
		return
	}

	reaction, err := sql.ParseUserReaction(reactionStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to parse reaction")
		return
	}

	err = h.commentService.AddOrUpdateOrRemoveReaction(int(user.ID), commentID, reaction)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to add or update or remove reaction")
		return
	}

	ref := r.Referer()
	if ref == "" {
		ref = "/"
	}

	comments.SetCommentRedirect(w, chapterID)
	http.Redirect(w, r, fmt.Sprintf("%s#comments", ref), http.StatusSeeOther)
}
