package novels

import (
	"fmt"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/config"
	"immodi/novel-site/pkg"
	"net/http"
	"time"
)

func IncrementNovelViews(service novels.NovelService, novelId int64) {
	err := service.IncrementNovelViewCount(novelId)
	if err != nil {
		fmt.Printf("Failed to increment novel views: %v\n", err)
	}
}

func IsAuthed(r *http.Request) bool {
	return r.Context().Value("user_id") != nil
}

func IsNovelBookMarked(r *http.Request, novelId int64, novelService novels.NovelService) bool {
	userID := IsAuthedUser(r)
	if userID == 0 {
		return false
	}

	isNovelBookMarked, err := novelService.IsNovelBookMarked(novelId, userID)
	if err != nil {
		return false
	}

	return isNovelBookMarked
}

func IsAuthedUser(r *http.Request) int64 {
	var tokenString string

	token, err := r.Cookie("auth_token")
	if err != nil {
		return 0
	}
	tokenString = token.Value

	userID, err := pkg.GetUserIDFromToken(tokenString)
	if err != nil {
		return 0
	}

	return userID
}

const (
	successCookieName = "successMessage"
	errorCookieName   = "errorMessage"
	flashDuration     = 60 * time.Second
)

// SetSuccessMessage sets a short-lived cookie with a success message.
func SetSuccessMessage(w http.ResponseWriter, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     successCookieName,
		Value:    msg,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.IsProduction,
		MaxAge:   int(flashDuration.Seconds()),
		Expires:  time.Now().Add(flashDuration),
	})
}

// SetErrorMessage sets a short-lived cookie with an error message.
func SetErrorMessage(w http.ResponseWriter, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     errorCookieName,
		Value:    msg,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(flashDuration.Seconds()),
		Expires:  time.Now().Add(flashDuration),
	})
}
