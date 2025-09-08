package search

import (
	"context"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/db/repositories"
	sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
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

func (s *searchService) GetGenreBySlug(slug string) (repositories.NovelGenre, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.NovelGenre, error) {
		return q.GetGenreBySlug(ctx, slug)
	})
}

func (s *searchService) GetTagBySlug(slug string) (repositories.NovelTag, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.NovelTag, error) {
		return q.GetTagBySlug(ctx, slug)
	})
}

func (s *searchService) GetAuthorNameBySlug(slug string) (string, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (string, error) {
		author, err := q.GetAuthorBySlug(ctx, slug)
		if err != nil {
			return pkg.SlugToTitle(slug), err
		}

		return author.Author, err
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

func (s *searchService) CountNovelsByGenre(genre string) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountNovelsByGenre(ctx, genre)
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

func (s *searchService) ListNovelsByAuthor(author string, offset int, limit int) (string, []repositories.Novel, error) {
	novels, err := db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNovelsByAuthorPaginated(ctx, repositories.ListNovelsByAuthorPaginatedParams{
			AuthorSlug: author,
			Offset:     int64(offset),
			Limit:      int64(limit),
		})
	})

	authorName, err := s.GetAuthorNameBySlug(author)
	return authorName, novels, err
}

func (s *searchService) ListNovelsByTag(tag string, offset int, limit int) (string, []repositories.Novel, error) {
	novels, err := db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNovelsByTagPaginated(ctx, repositories.ListNovelsByTagPaginatedParams{
			TagSlug: tag,
			Offset:  int64(offset),
			Limit:   int64(limit),
		})
	})

	tagName, err := s.GetTagBySlug(tag)
	return tagName.Tag, novels, err
}

func (s *searchService) ListNovelsByGenre(genre string, offset int, limit int) (string, []repositories.Novel, error) {
	novels, err := db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListNovelsByGenrePaginated(ctx, repositories.ListNovelsByGenrePaginatedParams{
			GenreSlug: genre,
			Offset:    int64(offset),
			Limit:     int64(limit),
		})
	})

	dbGenre, err := s.GetGenreBySlug(genre)
	return dbGenre.Genre, novels, err
}
