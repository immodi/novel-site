package index

import (
	"context"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/db/repositories"
)

type homeService struct {
	db *db.DBService
}

func NewHomeService(db *db.DBService) HomeService {
	return &homeService{db: db}
}

func (s *homeService) ListNewestNovels() ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNewestHomeNovels(ctx)
	})
}

func (s *homeService) ListHotNovels() ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListHotNovels(ctx)
	})
}

func (s *homeService) ListCompletedNovels() ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListCompletedNovels(ctx)
	})
}

func (s *homeService) ListGenres() ([]string, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		return q.ListAllGenres(ctx)
	})
}

func (s *homeService) GetLatestChapterByNovel(novelID int64) (repositories.Chapter, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.GetLatestChapterByNovel(ctx, novelID)
	})
}

func (s *homeService) ListGenresByNovel(novelID int64) ([]string, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		return q.ListGenresByNovel(ctx, novelID)
	})
}

func (s *homeService) CountChaptersByNovel(novelID int64) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountChaptersByNovel(ctx, novelID)
	})
}
