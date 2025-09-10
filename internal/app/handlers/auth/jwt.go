package auth

import (
	"fmt"
	"immodi/novel-site/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.JWTSecret)
var DefaultJwtDuration = 24 * time.Hour

func GenerateToken(userID int64, role string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses a JWT and returns the claims
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}

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
