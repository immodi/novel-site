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

func (s *novelService) GetNovelByName(name string) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByNameLike(ctx, name)
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

func (s *novelService) GetGenres(novelID int64) ([]string, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
		return q.ListGenresByNovel(ctx, novelID)
	})
}

func (s *novelService) GetTags(novelID int64) ([]string, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]string, error) {
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
			NovelID: novelID,
			Genre:   genre,
		})
	})
	return err
}

func (s *novelService) AddTagToNovel(novelID int64, tag string) error {
	_, err := db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (any, error) {
		return nil, q.AddTagToNovel(ctx, repositories.AddTagToNovelParams{
			NovelID: novelID,
			Tag:     tag,
		})
	})
	return err
}

func (s *novelService) CreateNovelWithDefaults(title string, isCompleted bool) (repositories.Novel, error) {
	minYear := 1980
	maxYear := time.Now().Year()
	randomYear := rand.Intn(maxYear-minYear+1) + minYear
	isCompletedInt := 0

	if isCompleted {
		isCompletedInt = 1
	}

	params := repositories.CreateNovelParams{
		Title:       title,
		Description: fmt.Sprintf("%s is a brand new story.", title),
		CoverImage:  "https://dummyimage.com/500x720/8a818a/ffffff",
		Author:      "default author",
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
