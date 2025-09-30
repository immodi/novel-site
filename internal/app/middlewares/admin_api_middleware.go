package middlewares

import (
	"context"
	"encoding/json"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
	"log"
	"net/http"
	"strings"
)

type apiError struct {
	Error string `json:"error"`
}

func ApiAdminRoleMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("missing Authorization header")
				writeJSONError(w, http.StatusUnauthorized, "missing Authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				log.Println("invalid Authorization header format")
				writeJSONError(w, http.StatusUnauthorized, "invalid Authorization header format")
				return
			}
			tokenString := parts[1]

			// Parse JWT
			claims, err := pkg.ParseToken(tokenString)
			if err != nil {
				log.Println("failed to parse token:", err)
				writeJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// Check if role is strictly admin
			role, ok := claims["role"].(string)
			if !ok || role != string(sql.UserRoleAdmin) {
				log.Println("role mismatch or missing")
				writeJSONError(w, http.StatusForbidden, "admin role required")
				return
			}

			// Store user_id in context
			userID := int64(claims["user_id"].(float64))
			ctx := context.WithValue(r.Context(), "user_id", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError{Error: msg})
}
