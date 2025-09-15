package load

import (
	"context"
	"immodi/novel-site/internal/app/services/chapters"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/db/repositories"
)

type loadService struct {
	db             *db.DBService
	novelService   novels.NovelService
	chapterService chapters.ChapterService
}

func NewLoadService(db *db.DBService, novelService novels.NovelService, chapterService chapters.ChapterService) LoadService {
	return &loadService{db: db, novelService: novelService, chapterService: chapterService}
}

func (s *loadService) CreateNovel(params repositories.CreateNovelParams) (repositories.Novel, error) {
	return s.novelService.CreateNovel(params)
}

func (s *loadService) GetNovelByExactName(name string) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByExactName(ctx, name)
	})
}

func (s *loadService) GetNovelById(id int64) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByID(ctx, id)
	})
}

func (s *loadService) GetLastNovelChapter(novelID int64) (repositories.Chapter, error) {
	return s.novelService.GetLastChapter(novelID)
}

func (s *loadService) AddBulkTagsToNovel(novelID int64, tags []string) error {
	return s.novelService.AddBulkTagsToNovel(novelID, tags)
}

func (s *loadService) AddBulkGenresToNovel(novelID int64, genres []string) error {
	return s.novelService.AddBulkGenresToNovel(novelID, genres)
}

func (s *loadService) DeleteNovel(novelID int64) error {
	return s.novelService.DeleteNovel(novelID)
}

func (s *loadService) CreateBulkChapters(chapters []repositories.CreateChapterParams) error {
	return s.chapterService.CreateBulkChapters(chapters)
}
