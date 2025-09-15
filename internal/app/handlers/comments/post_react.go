package comments

import (
	"fmt"
	"immodi/novel-site/internal/app/handlers"
	sql "immodi/novel-site/internal/db/schema"
	"log"
	"net/http"
	"strconv"
)

func (h *CommentHandler) PostReact(w http.ResponseWriter, r *http.Request) {
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

	novelIDStr := r.FormValue("novelID")
	commentIDStr := r.FormValue("commentID")
	reactionStr := r.FormValue("reaction")

	novelID, err := strconv.Atoi(novelIDStr)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to parse novelID")
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

	SetCommentRedirect(w, novelID)
	http.Redirect(w, r, fmt.Sprintf("%s#comments", ref), http.StatusSeeOther)
}
