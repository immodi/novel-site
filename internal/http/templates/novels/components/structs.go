package components

type Novel struct {
	Name                string
	Description         string
	Author              string
	Genres              []string
	Tags                []string
	Status              string
	CoverImage          string
	CurrentPage         int
	TotalPages          int
	TotalChaptersNumber int
	Chapters            []Chapter
	LastChapterName     string
	LastUpdated         string
}

type Chapter struct {
	Title  string
	Number int
}
