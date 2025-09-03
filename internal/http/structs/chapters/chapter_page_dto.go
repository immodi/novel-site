package chaptersdtostructs

type ChapterPage struct {
	NovelName      string
	ChapterTitle   string
	ChapterContent string
	PrevChapter    *int
	NextChapter    *int
}
