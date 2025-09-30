package index

import "immodi/novel-site/internal/db/repositories"

type HomeService interface {
	ListAllNovels() ([]repositories.Novel, error)

	ListNewestNovels() ([]repositories.Novel, error)
	ListHotNovels() ([]repositories.Novel, error)
	ListCompletedNovels() ([]repositories.Novel, error)
	ListGenresByNovel(novelID int64) ([]repositories.NovelGenre, error)
	ListGenres() ([]string, error)

	GetLatestChapterByNovel(novelID int64) (repositories.Chapter, error)

	CountChaptersByNovel(novelID int64) (int64, error)
}
