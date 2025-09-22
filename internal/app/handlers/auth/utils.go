package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

var CookieName = "auth_token"
var OauthStateName = "oauth_token"

func removeAuthCookie(w http.ResponseWriter, isProduction bool) {
	var secure bool = false
	if isProduction {
		secure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		MaxAge:   -1,
	})
}

func setAuthCookie(w http.ResponseWriter, token string, duration time.Duration, isProduction bool) {
	var secure bool = false

	if isProduction {
		secure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		MaxAge:   int(duration.Seconds()),
	})
}

func authenticateRegisterRequest(username, email, password, confirmPassword, terms string) []string {
	var errors []string

	if username == "" {
		errors = append(errors, "Username is required")
	}
	if email == "" {
		errors = append(errors, "Email is required")
	}
	if password == "" {
		errors = append(errors, "Password is required")
	}

	if confirmPassword == "" {
		errors = append(errors, "Please confirm your password")
	}

	if password != confirmPassword {
		errors = append(errors, "Passwords do not match")
	}

	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}

	if terms != "on" {
		errors = append(errors, "You must accept the terms and privacy policy")
	}

	return errors
}

func authenticateLoginRequest(email, password string) []string {
	var errors []string

	if email == "" {
		errors = append(errors, "Email is required")
	}
	if password == "" {
		errors = append(errors, "Password is required")
	}

	return errors
}

func GenerateStateOauthCookie(w http.ResponseWriter) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Println("Failed to generate random state:", err)
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     OauthStateName,
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   300,
	})
	return state, nil
}
