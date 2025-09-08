package searchresutlsdto

type SearchResultDto struct {
	ID                   int
	Name                 string
	Slug                 string
	CoverImage           string
	Author               string
	Status               string // e.g., "Ongoing", "Completed"
	LastestChapterName   string
	LastestChapterNumber int
	LastUpdated          string
}
