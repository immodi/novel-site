package pkg

import "net/http"

// GetAndClearCookie retrieves a cookie value and deletes it immediately.
func GetAndClearCookie(w http.ResponseWriter, r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}

	// Clear it (flash behavior)
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	return cookie.Value
}
