package comments

import (
	"immodi/novel-site/internal/app/handlers/auth"
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

func IsAuthedUser(r *http.Request) int64 {
	var tokenString string

	token, err := r.Cookie("auth_token")
	if err != nil {
		return 0
	}
	tokenString = token.Value

	userID, err := auth.GetUserIDFromToken(tokenString)
	if err != nil {
		return 0
	}

	return userID
}
