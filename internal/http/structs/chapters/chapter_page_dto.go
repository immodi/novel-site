package chaptersdtostructs

type ChapterPage struct {
	NovelName      string
	NovelSlug      string
	ChapterTitle   string
	ChapterContent string
	ChapterID      int
	ChapterNumber  int
	PrevChapter    *int
	NextChapter    *int
}
