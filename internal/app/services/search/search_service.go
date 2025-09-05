package search

import (
	"context"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/db/repositories"
	sql "immodi/novel-site/internal/db/schema"
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

func (s *searchService) CountSortedNovels(collection sql.Collection) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		if collection == sql.CollectionCompleted {
			return q.CountCompletedNovels(ctx)
		}

		if collection == sql.CollectionOnGoing {
			return q.CountOnGoingNovels(ctx)
		}

		return q.CountNovels(ctx)
	})
}

func (s *searchService) ListSortedNovels(collection sql.Collection, offset int, limit int) ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		switch collection {
		case sql.CollectionCompleted:
			return q.ListCompletedNovelsPaginated(ctx, repositories.ListCompletedNovelsPaginatedParams{
				Offset: int64(offset),
				Limit:  int64(limit),
			})
		case sql.CollectionHot:
			return q.ListHotNovelsPaginated(ctx, repositories.ListHotNovelsPaginatedParams{
				Offset: int64(offset),
				Limit:  int64(limit),
			})
		case sql.CollectionLatest:
			return q.ListNewestHomeNovelsPaginated(ctx, repositories.ListNewestHomeNovelsPaginatedParams{
				Offset: int64(offset),
				Limit:  int64(limit),
			})
		case sql.CollectionOnGoing:
			return q.ListOnGoingNovelsPaginated(ctx, repositories.ListOnGoingNovelsPaginatedParams{
				Offset: int64(offset),
				Limit:  int64(limit),
			})
		default:
			return q.ListHotNovelsPaginated(ctx, repositories.ListHotNovelsPaginatedParams{
				Offset: int64(offset),
				Limit:  int64(limit),
			})
		}

	})
}

func (s *searchService) CountNovelsByGenre(genre sql.Genre) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountNovelsByGenre(ctx, string(genre))
	})
}

func (s *searchService) ListNovelsByGenre(genre sql.Genre, offset int, limit int) ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNovelsByGenrePaginated(ctx, repositories.ListNovelsByGenrePaginatedParams{
			Genre:  string(genre),
			Offset: int64(offset),
			Limit:  int64(limit),
		})
	})
}

func (s *searchService) CountNovelsByAuthor(author string) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountNovelsByAuthor(ctx, author)
	})
}

func (s *searchService) CountNovelsByTag(tag string) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountNovelsByTag(ctx, tag)
	})
}

func (s *searchService) ListNovelsByAuthor(author string, offset int, limit int) ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNovelsByAuthorPaginated(ctx, repositories.ListNovelsByAuthorPaginatedParams{
			Author: author,
			Offset: int64(offset),
			Limit:  int64(limit),
		})
	})
}

func (s *searchService) ListNovelsByTag(tag string, offset int, limit int) ([]repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNovelsByTagPaginated(ctx, repositories.ListNovelsByTagPaginatedParams{
			Tag:    tag,
			Offset: int64(offset),
			Limit:  int64(limit),
		})
	})
}
