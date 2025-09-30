package sitemap

import (
	"context"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/db/repositories"
)

type sitemapService struct {
	db *db.DBService
}

func NewSitemapService(db *db.DBService) SitemapService {
	return &sitemapService{db: db}
}

func (s *sitemapService) GetAllNovels() ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListAllNovels(ctx)
	})
}

func (s *sitemapService) GetAllGenres() ([]string, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		return q.ListAllGenreSlugs(ctx)
	})
}

func (s *sitemapService) GetAllTags() ([]string, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		return q.ListAllTagSlugs(ctx)
	})
}
