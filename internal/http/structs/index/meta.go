package indexdtostructs

const DOMAIN = "https://inovelhub.com"
const SITE_NAME = "INovelHub"

type MetaDataStruct struct {
	IsRendering       bool
	Title             string
	Description       string
	Keywords          string
	OgURL             string
	Canonical         string
	CoverImage        string
	Genres            []string
	Author            string
	AuthorLink        string
	NovelName         string
	ReadURL           string
	Status            string
	UpdateTime        string
	LatestChapterName string
	LatestChapterURL  string
}
