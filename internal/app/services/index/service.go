package index

import "immodi/novel-site/internal/db/repositories"

type HomeService interface {
	ListNewestNovels() ([]repositories.Novel, error)
	ListHotNovels() ([]repositories.Novel, error)
	ListCompletedNovels() ([]repositories.Novel, error)
	GetLatestChapterByNovel(novelID int64) (repositories.Chapter, error)
	ListGenresByNovel(novelID int64) ([]string, error)
	CountChaptersByNovel(novelID int64) (int64, error)
	ListGenres() ([]string, error)
}
