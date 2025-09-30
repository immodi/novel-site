package admin

import (
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/payloads/admin"
)

func DbNovelsToAdminPanelNovelsMapper(dbNovels []repositories.Novel) []admin.Novel {
	var adminPanelNovels []admin.Novel

	for _, novel := range dbNovels {
		var novelStatus string = "Completed"
		if novel.IsCompleted == 1 {
			novelStatus = "Ongoing"
		}

		adminPanelNovels = append(adminPanelNovels, admin.Novel{
			ID:     novel.ID,
			Title:  novel.Title,
			Views:  novel.ViewCount,
			Author: novel.Author,
			Status: novelStatus,
		})
	}

	return adminPanelNovels
}
