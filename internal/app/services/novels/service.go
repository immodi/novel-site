package novels

import (
	"immodi/novel-site/internal/db/repositories"
)

type NovelService interface {
	GetNovelBySlug(slug string) (repositories.Novel, error)
	GetChapters(novelID, page int) ([]repositories.Chapter, error)
	GetChapterByID(chapterID int64) (repositories.Chapter, error)

	GetGenres(novelID int64) ([]repositories.NovelGenre, error)
	GetAllGenres() ([]string, error)

	GetTags(novelID int64) ([]repositories.NovelTag, error)
	FilterTagsByName(tag string) ([]string, error)

	CountChapters(novelID int64) (int, error)
	GetLastChapter(novelID int64) (repositories.Chapter, error)
	CreateNovel(params repositories.CreateNovelParams) (repositories.Novel, error)
	CreateNovelWithDefaults(title string, isCompleted bool) (repositories.Novel, error)

	AddBulkGenresToNovel(novelID int64, genres []string) error
	AddBulkTagsToNovel(novelID int64, tags []string) error

	AddGenreToNovel(novelID int64, genre string) error
	AddTagToNovel(novelID int64, tag string) error

	IncrementNovelViewCount(novelID int64) error
	DeleteNovel(novelID int64) error

	IsNovelBookMarked(novelID int64, userID int64) (bool, error)
	GetLastReadChapterID(userID, novelID int64) (int64, error)
}
