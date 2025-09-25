package middlewares

import (
	"context"
	"immodi/novel-site/internal/app/handlers/auth"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
)

func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(auth.CookieName)
			if err != nil {
				log.Println("no identity cookie")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			tokenString := cookie.Value
			if tokenString == "" {
				log.Println("empty identity cookie")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Parse JWT
			claims, err := pkg.ParseToken(tokenString)
			if err != nil {
				log.Println("failed to parse identity token")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Check role
			role, ok := claims["role"].(string)
			if !ok {
				log.Println("failed to get role from identity token")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			// allow if role matches OR role is admin
			if role != requiredRole && role != string(sql.UserRoleAdmin) {
				log.Println("role mismatch")
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
