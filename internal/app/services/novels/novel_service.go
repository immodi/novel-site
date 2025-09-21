package novels

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/pkg"
	"math/rand"
	"time"
)

type novelService struct {
	db *db.DBService
}

func New(db *db.DBService) NovelService {
	return &novelService{db: db}
}

func (s *novelService) GetNovelBySlug(slug string) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelBySlug(ctx, slug)
	})
}

func (s *novelService) GetChapters(novelID, page int) ([]repositories.Chapter, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Chapter, error) {
		return q.ListChaptersByNovelPaginated(ctx, repositories.ListChaptersByNovelPaginatedParams{
			NovelID: int64(novelID),
			Limit:   pkg.PAGE_LIMIT,
			Offset:  int64(pkg.PAGE_LIMIT * (page - 1)),
		})
	})
}

func (s *novelService) CountChapters(novelID int64) (int, error) {
	total, err := db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountChaptersByNovel(ctx, novelID)
	})
	return int(total), err
}

func (s *novelService) GetLastChapter(novelID int64) (repositories.Chapter, error) {
	total, err := s.CountChapters(novelID)
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

func (s *novelService) GetGenres(novelID int64) ([]repositories.NovelGenre, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.NovelGenre, error) {
		return q.ListGenresByNovel(ctx, novelID)
	})
}

func (s *novelService) GetTags(novelID int64) ([]repositories.NovelTag, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.NovelTag, error) {
		return q.GetNovelTags(ctx, novelID)
	})
}

func (s *novelService) CreateNovel(params repositories.CreateNovelParams) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.CreateNovel(ctx, params)
	})
}

func (s *novelService) AddGenreToNovel(novelID int64, genre string) error {
	_, err := db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (any, error) {
		return nil, q.AddGenreToNovel(ctx, repositories.AddGenreToNovelParams{
			NovelID:   novelID,
			Genre:     genre,
			GenreSlug: pkg.TitleToSlug(genre),
		})
	})
	return err
}

func (s *novelService) AddTagToNovel(novelID int64, tag string) error {
	_, err := db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (any, error) {
		return nil, q.AddTagToNovel(ctx, repositories.AddTagToNovelParams{
			NovelID: novelID,
			Tag:     tag,
			TagSlug: pkg.TitleToSlug(tag),
		})
	})
	return err
}

func (s *novelService) AddBulkGenresToNovel(novelID int64, genres []string) error {
	return db.ExecuteTx(s.db, func(ctx context.Context, q *repositories.Queries) error {
		for _, genre := range genres {
			if genre == "" {
				continue
			}
			err := q.AddGenreToNovel(ctx, repositories.AddGenreToNovelParams{
				NovelID:   novelID,
				Genre:     genre,
				GenreSlug: pkg.TitleToSlug(genre),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *novelService) AddBulkTagsToNovel(novelID int64, tags []string) error {
	return db.ExecuteTx(s.db, func(ctx context.Context, q *repositories.Queries) error {
		for _, tag := range tags {
			if tag == "" {
				continue
			}
			err := q.AddTagToNovel(ctx, repositories.AddTagToNovelParams{
				NovelID: novelID,
				Tag:     tag,
				TagSlug: pkg.TitleToSlug(tag),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *novelService) CreateNovelWithDefaults(title string, isCompleted bool) (repositories.Novel, error) {
	minYear := 1980
	maxYear := time.Now().Year()
	randomYear := rand.Intn(maxYear-minYear+1) + minYear
	isCompletedInt := 0
	author := "Default Author"
	if isCompleted {
		isCompletedInt = 1
	}

	params := repositories.CreateNovelParams{
		Title:       title,
		Slug:        pkg.TitleToSlug(title),
		Description: fmt.Sprintf("%s is a brand new story.", title),
		CoverImage:  "https://dummyimage.com/500x720/8a818a/ffffff",
		Author:      author,
		AuthorSlug:  pkg.TitleToSlug(author),
		Publisher:   "Default Publisher",
		ReleaseYear: int64(pkg.GetRandomYear(randomYear)),
		IsCompleted: int64(isCompletedInt),
		UpdateTime:  pkg.GetCurrentTimeRFC3339(),
	}

	return s.CreateNovel(params)
}

func (s *novelService) IncrementNovelViewCount(novelID int64) error {
	return db.Execute(s.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.IncrementNovelViewCount(ctx, novelID)
	})
}

func (s *novelService) DeleteNovel(novelID int64) error {
	return db.Execute(s.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.DeleteNovel(ctx, novelID)
	})
}

func (s *novelService) IsNovelBookMarked(novelID int64, userID int64) (bool, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (bool, error) {
		value, err := q.IsNovelBookmarked(ctx, repositories.IsNovelBookmarkedParams{
			NovelID: novelID,
			UserID:  userID,
		})

		if err != nil {
			return false, err
		}

		return value == 1, nil
	})
}

func (s *novelService) GetLastReadChapterID(userID, novelID int64) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		chapterID, err := q.GetLastReadChapterID(ctx, repositories.GetLastReadChapterIDParams{
			UserID:  userID,
			NovelID: novelID,
		})

		if err != nil || !chapterID.Valid {
			return 0, err
		}

		return chapterID.Int64, nil
	})
}

func (s *novelService) GetChapterByID(chapterID int64) (repositories.Chapter, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.GetChapterByID(ctx, chapterID)
	})
}
