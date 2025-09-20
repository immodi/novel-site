package chapters

import "immodi/novel-site/internal/db/repositories"

type ChapterService interface {
	GetNovelBySlug(slug string) (repositories.Novel, error)
	CountChaptersByNovel(novelID int64) (int64, error)

	GetChapterByNumber(novelID int64, chapterNumber int64) (repositories.Chapter, error)
	ListChaptersByNovel(novelID int64) ([]repositories.Chapter, error)

	CreateChapterWithDefaults(novelID int64) (repositories.Chapter, error)
	CreateBulkChapters(chapters []repositories.CreateChapterParams) error

	IncrementNovelViewCount(novelID int64) error
}
