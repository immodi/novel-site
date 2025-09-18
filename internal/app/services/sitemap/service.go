package sitemap

import "immodi/novel-site/internal/db/repositories"

type SitemapService interface {
	GetAllTags() ([]string, error)
	GetAllGenres() ([]string, error)
	GetAllNovels() ([]repositories.Novel, error)
}
