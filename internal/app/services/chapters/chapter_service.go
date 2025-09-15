package chapters

import (
	"context"
	"fmt"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/pkg"
)

type chapterService struct {
	db           *db.DBService
	novelService novels.NovelService
}

func New(db *db.DBService, novelService novels.NovelService) ChapterService {
	return &chapterService{db: db, novelService: novelService}
}

func (s *chapterService) GetNovelBySlug(slug string) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelBySlug(ctx, slug)
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
			ReleaseDate:   pkg.GetCurrentTimeRFC3339(),
		})
	})
}

func (s *chapterService) CreateBulkChapters(chapters []repositories.CreateChapterParams) error {
	return db.ExecuteTx(s.db, func(ctx context.Context, q *repositories.Queries) error {
		if len(chapters) == 0 {
			return nil
		}

		novelID := chapters[0].NovelID

		for _, chapter := range chapters {
			if _, err := q.CreateChapter(ctx, chapter); err != nil {
				return err
			}
		}

		_, err := q.UpdateNovelPartial(ctx, repositories.UpdateNovelPartialParams{
			ID:         novelID,
			UpdateTime: pkg.GetCurrentTimeRFC3339(),
		})
		return err
	})
}

func (s *chapterService) ListChaptersByNovel(novelID int64) ([]repositories.Chapter, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Chapter, error) {
		return q.ListChaptersByNovel(ctx, novelID)
	})
}
