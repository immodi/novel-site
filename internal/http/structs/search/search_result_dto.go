package searchresutlsdto

type SearchResultDto struct {
	ID                   int
	Name                 string
	CoverImage           string
	Author               string
	Status               string // e.g., "Ongoing", "Completed"
	LastestChapterName   string
	LastestChapterNumber int
	LastUpdated          string
	Genres               []string
	ChaptersCount        int
	Description          string
}
