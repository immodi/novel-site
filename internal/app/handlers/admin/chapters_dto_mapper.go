package admin

import (
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/internal/http/payloads/admin"
)

func DbChaptersToAdminPanelChaptersMapper(dbChapters []repositories.Chapter) []admin.Chapter {
	var adminPanelChapters []admin.Chapter

	for _, chapter := range dbChapters {
		adminPanelChapters = append(adminPanelChapters, admin.Chapter{
			ID:          chapter.ID,
			NovelID:     chapter.NovelID,
			ReleaseDate: chapter.ReleaseDate,
			Title:       chapter.Title,
			Content:     chapter.Content,
		})
	}

	return adminPanelChapters
}
