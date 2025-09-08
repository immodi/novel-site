package load

import "immodi/novel-site/internal/db/repositories"

type LoadService interface {
	GetNovelByExactName(name string) (repositories.Novel, error)
	CreateNovel(params repositories.CreateNovelParams) (repositories.Novel, error)
	AddGenreToNovel(novelID int64, genre string) error
	AddTagToNovel(novelID int64, tag string) error
	DeleteNovel(novelID int64) error
}
