package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"immodi/novel-site/internal/app/handlers"
	"immodi/novel-site/internal/config"
	"immodi/novel-site/internal/db/repositories"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

var googleApiEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo"
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  fmt.Sprintf("%s/auth/google/callback", config.SiteURL),
	ClientID:     config.GoogleClientID,
	ClientSecret: config.GoogleClientSecret,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func (h *AuthHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	oauthState, err := GenerateStateOauthCookie(w)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to generate state")
		return
	}

	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Retrieve state from cookie
	cookie, err := r.Cookie(OauthStateName)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to get cookie")
		return
	}
	oauthState := cookie.Value

	if r.URL.Query().Get("state") != oauthState {
		handlers.NotFoundHandler(w, r)
		log.Println("invalid oauth state")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.NotFoundHandler(w, r)
		log.Println("failed to get code")
		return
	}

	// Exchange code for token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to exchange code for token")
		return
	}

	// Get user info
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get(googleApiEndpoint)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to get user info")
		return
	}
	defer resp.Body.Close()

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to decode user info")
		return
	}

	// Check if user exists
	_, exists, err := h.authService.GetUserByEmail(user.Email)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to get user")
		return
	}

	if !exists {
		// Create user
		userPasswordHash, err := pkg.HashPassword(user.ID)
		if err != nil {
			handlers.ServerErrorHandler(w, r)
			log.Println(err.Error(), "failed to hash password")
			return
		}

		dbUserID, err := h.authService.RegisterUser(repositories.CreateUserParams{
			Username:     user.Name,
			Email:        user.Email,
			PasswordHash: userPasswordHash,
			CreatedAt:    pkg.GetCurrentTimeRFC3339(),
			Role:         string(sql.UserRoleUser),
		})

		if err != nil {
			handlers.ServerErrorHandler(w, r)
			log.Println(err.Error(), "failed to create user")
			return
		}

		err = h.authService.UpdateUserImage(dbUserID, user.Picture)
		if err != nil {
			handlers.ServerErrorHandler(w, r)
			log.Println(err.Error(), "failed to update user image")
			return
		}
	}

	dbUser, _, err := h.authService.GetUserByEmail(user.Email)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to login user")
		return
	}

	// Generate JWT token
	jwtToken, err := pkg.GenerateToken(dbUser.ID, dbUser.Role, pkg.DefaultJwtDuration)
	if err != nil {
		handlers.ServerErrorHandler(w, r)
		log.Println(err.Error(), "failed to generate token")
		return
	}

	// Set JWT cookie
	setAuthCookie(w, jwtToken, pkg.DefaultJwtDuration, config.IsProduction)
	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}
