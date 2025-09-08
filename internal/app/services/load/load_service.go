package load

import (
	"context"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/db/repositories"
)

type loadService struct {
	db           *db.DBService
	novelService novels.NovelService
}

func NewLoadService(db *db.DBService, novelService novels.NovelService) LoadService {
	return &loadService{db: db, novelService: novelService}
}

func (s *loadService) CreateNovel(params repositories.CreateNovelParams) (repositories.Novel, error) {
	return s.novelService.CreateNovel(params)
}

func (s *loadService) GetNovelByExactName(name string) (repositories.Novel, error) {
	return db.ExecuteWithResult(s.db, func(ctx context.Context, q *repositories.Queries) (repositories.Novel, error) {
		return q.GetNovelByExactName(ctx, name)
	})
}

func (s *loadService) AddGenreToNovel(novelID int64, genre string) error {
	return s.novelService.AddGenreToNovel(novelID, genre)
}

func (s *loadService) AddTagToNovel(novelID int64, tag string) error {
	return s.novelService.AddTagToNovel(novelID, tag)
}

func (s *loadService) DeleteNovel(novelID int64) error {
	return s.novelService.DeleteNovel(novelID)
}
