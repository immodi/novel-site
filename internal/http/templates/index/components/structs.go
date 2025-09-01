package components

type HomeNovelDto struct {
	Name                 string
	CoverImage           string
	LastestChapterNumber int
	LastestChapterName   string
	LastUpdated          string
	Status               string
	ChaptersCount        int
	Genres               []string
}
