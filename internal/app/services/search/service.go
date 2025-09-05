package search

import (
	"immodi/novel-site/internal/db/repositories"
	sql "immodi/novel-site/internal/db/schema"
)

type SearchService interface {
	SearchNovelsPaginated(name string, offset, limit int) ([]repositories.Novel, error)
	CountTotalSearchedNovels(name string) (int64, error)
	GetLastChapter(novelID int64) (repositories.Chapter, error)
	CountSortedNovels(collection sql.Collection) (int64, error)
	ListSortedNovels(collection sql.Collection, offset, limit int) ([]repositories.Novel, error)
	CountNovelsByGenre(genre sql.Genre) (int64, error)
	ListNovelsByGenre(genre sql.Genre, offset, limit int) ([]repositories.Novel, error)
	CountNovelsByTag(tag string) (int64, error)
	ListNovelsByTag(tag string, offset, limit int) ([]repositories.Novel, error)
	CountNovelsByAuthor(author string) (int64, error)
	ListNovelsByAuthor(author string, offset, limit int) ([]repositories.Novel, error)
}
