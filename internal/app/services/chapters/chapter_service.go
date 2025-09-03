package chapters

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/pkg"
)

type chapterService struct {
	db *db.DBService
}

func New(db *db.DBService) ChapterService {
	return &chapterService{db: db}
}

func (s *chapterService) GetNovelByNameLike(name string) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByNameLike(ctx, name)
	})
}

func (s *chapterService) CountChaptersByNovel(novelID int64) (int64, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountChaptersByNovel(ctx, novelID)
	})
}

func (s *chapterService) GetChapterByNumber(novelID int64, chapterNumber int64) (repositories.Chapter, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		return q.GetChapterByNumber(ctx, repositories.GetChapterByNumberParams{
			NovelID:       novelID,
			ChapterNumber: chapterNumber,
		})
	})
}

func (s *chapterService) CreateChapterWithDefaults(novelID int64) (repositories.Chapter, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Chapter, error) {
		// First, count chapters
		count, err := q.CountChaptersByNovel(ctx, novelID)
		if err != nil {
			return repositories.Chapter{}, err
		}

		_, err = db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
			return q.UpdateNovelPartial(ctx, repositories.UpdateNovelPartialParams{
				ID:         novelID,
				UpdateTime: pkg.GetCurrentTimeRFC3339(),
			})
		})

		if err != nil {
			return repositories.Chapter{}, err
		}

		// Then create chapter
		return q.CreateChapter(ctx, repositories.CreateChapterParams{
			NovelID:       novelID,
			Title:         fmt.Sprintf("Chapter %d", count+1),
			ChapterNumber: count + 1,
			Content:       pkg.LoremText(40),
		})
	})
}

func (s *chapterService) ListChaptersByNovel(novelID int64) ([]repositories.Chapter, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Chapter, error) {
		return q.ListChaptersByNovel(ctx, novelID)
	})
}
