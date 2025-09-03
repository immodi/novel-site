package novelsdtostructs

type Novel struct {
	Name                string
	Description         string
	Author              string
	Genres              []string
	Tags                []string
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
}
