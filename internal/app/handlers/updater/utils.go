package updater

import (
	"fmt"
	"immodi/novel-site/internal/config"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		adminOrigin := strings.TrimSuffix(config.AdminSiteURL, "/")
		requestOrigin := strings.TrimSuffix(r.Header.Get("Origin"), "/")

		return requestOrigin == adminOrigin
	},
}

func (h *UpdaterHandler) Validate(r *http.Request) error {
	cookie, err := r.Cookie("admin_auth_cookie")
	if err != nil {
		return fmt.Errorf("missing auth cookie")
	}

	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return fmt.Errorf("failed to decode cookie")
	}

	token := strings.Trim(decodedValue, `"`)
	if token == "" {
		return fmt.Errorf("missing token in cookie")
	}

	userID, err := pkg.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	user, err := h.profileService.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("could not get the user from the token")
	}

	if user.Role != string(sql.UserRoleAdmin) {
		return fmt.Errorf("user is not an admin")
	}

	return nil
}
