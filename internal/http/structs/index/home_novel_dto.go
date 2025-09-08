package indexdtostructs

import "immodi/novel-site/internal/db/repositories"

type HomeNovelDto struct {
	Name                 string
	Slug                 string
	CoverImage           string
	LastestChapterNumber int
	LastestChapterName   string
	LastUpdated          string
	Status               string
	ChaptersCount        int
	Genres               []repositories.NovelGenre
}
