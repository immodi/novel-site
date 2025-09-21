package pkg

import (
	"fmt"
	"net/http"
)

// GetUserIDFromToken parses the token and extracts the user_id claim as int64
func GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	// Ensure user_id exists and is a number
	uid, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid or missing user_id in token")
	}

	return int64(uid), nil
}

func IsAuthedUser(r *http.Request) int64 {
	var tokenString string

	token, err := r.Cookie("auth_token")
	if err != nil {
		return 0
	}
	tokenString = token.Value

	userID, err := GetUserIDFromToken(tokenString)
	if err != nil {
		return 0
	}

	return userID
}
