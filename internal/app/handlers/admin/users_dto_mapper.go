package admin

import (
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/payloads/admin"
)

func DbUsersToAdminPanelUsersMapper(dbUsers []repositories.User) []admin.User {
	var adminPanelUsers []admin.User

	for _, user := range dbUsers {
		adminPanelUsers = append(adminPanelUsers, admin.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			Image:     user.Image,
		})
	}

	return adminPanelUsers
}
