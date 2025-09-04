package search

import (
	"context"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/db/repositories"
	"strings"
)

type searchService struct {
	db          *db.DBService
	homeService index.HomeService
}

func NewSearchService(db *db.DBService, homeService index.HomeService) SearchService {
	return &searchService{db: db, homeService: homeService}
}

func (s *searchService) GetLastChapter(novelID int64) (repositories.Chapter, error) {
	total, err := s.homeService.CountChaptersByNovel(novelID)
	if err != nil {
		return repositories.Chapter{}, err
	}

	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.GetChapterByNumber(ctx, repositories.GetChapterByNumberParams{
			NovelID:       novelID,
			ChapterNumber: int64(total),
		})
	})
}

func (s *searchService) SearchNovelsPaginated(name string, offset, limit int) ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.SearchNovels(ctx, repositories.SearchNovelsParams{
			Search: strings.ToLower(name),
			Offset: int64(offset),
			Limit:  int64(limit),
		})
	})
}

func (s *searchService) CountTotalSearchedNovels(name string) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountSearchNovels(ctx, name)
	})
}
