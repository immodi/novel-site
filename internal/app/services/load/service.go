package load

import "immodi/novel-site/internal/db/repositories"

type LoadService interface {
	GetNovelByExactName(name string) (repositories.Novel, error)
	GetNovelById(id int64) (repositories.Novel, error)

	GetLastNovelChapter(novelID int64) (repositories.Chapter, error)

	CreateNovel(params repositories.CreateNovelParams) (repositories.Novel, error)
	CreateBulkChapters(chapters []repositories.CreateChapterParams) error

	AddBulkGenresToNovel(novelID int64, genres []string) error

	AddBulkTagsToNovel(novelID int64, tags []string) error
	DeleteNovel(novelID int64) error
}
