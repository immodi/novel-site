package comments

import (
	"immodi/novel-site/internal/config"
	"net/http"
	"strconv"
	"time"
)

const CommentRedirectCookie = "comment_redirect"

var commentRedirectDuration = 5 * time.Second // short-lived

func SetCommentRedirect(w http.ResponseWriter, id int) {
	http.SetCookie(w, &http.Cookie{
		Name:     CommentRedirectCookie,
		Value:    strconv.Itoa(id),
		Path:     "/",
		HttpOnly: true,
		Secure:   config.IsProduction,
		MaxAge:   int(commentRedirectDuration.Seconds()),
		Expires:  time.Now().Add(commentRedirectDuration),
		SameSite: http.SameSiteDefaultMode,
	})
}
