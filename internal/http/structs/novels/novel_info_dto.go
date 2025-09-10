package novelsdtostructs

import "immodi/novel-site/internal/db/repositories"

type Novel struct {
	Name                string
	Slug                string
	Description         string
	Author              string
	AuthorSlug          string
	Genres              []repositories.NovelGenre
	Tags                []repositories.NovelTag
	Views               string
	Status              string
	CoverImage          string
	Publisher           string
	ReleaseYear         int
	CurrentPage         int
	TotalPages          int
	TotalChaptersNumber int
	Chapters            []Chapter
	LastChapterName     string
	LastUpdated         string
	IsNovelBookMarked   bool
	SuccessMessage      string
	ErrorMessage        string
}
