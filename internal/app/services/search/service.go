package search

import "immodi/novel-site/internal/db/repositories"

type SearchService interface {
	SearchNovelsPaginated(name string, offset, limit int) ([]repositories.Novel, error)
	CountTotalSearchedNovels(name string) (int64, error)
	GetLastChapter(novelID int64) (repositories.Chapter, error)
}
