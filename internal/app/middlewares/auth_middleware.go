package middlewares

import (
	"context"
	"immodi/novel-site/internal/app/handlers/auth"
	"net/http"
)

func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(auth.CookieName)
			if err != nil {
				// Cookie missing
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			tokenString := cookie.Value
			if tokenString == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Parse JWT
			claims, err := auth.ParseToken(tokenString)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Check role
			role, ok := claims["role"].(string)
			if !ok || role != requiredRole {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			// Store user_id in context
			userID := int64(claims["user_id"].(float64))
			ctx := context.WithValue(r.Context(), "user_id", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
