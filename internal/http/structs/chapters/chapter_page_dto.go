package chaptersdtostructs

type ChapterPage struct {
	NovelName      string
	NovelSlug      string
	ChapterTitle   string
	ChapterContent string
	PrevChapter    *int
	NextChapter    *int
}
