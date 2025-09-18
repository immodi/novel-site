package chapters

import (
	"immodi/novel-site/internal/db/repositories"
	chaptersdtostructs "immodi/novel-site/internal/http/structs/chapters"
	novelsdtostructs "immodi/novel-site/internal/http/structs/novels"
)

func DbChaptersToChaptersMapper(dbChapters []repositories.Chapter) []novelsdtostructs.Chapter {
	chapters := make([]novelsdtostructs.Chapter, len(dbChapters))

	for i, dbChapter := range dbChapters {
		chapters[i] = novelsdtostructs.Chapter{
			Title:  dbChapter.Title,
			Number: int(dbChapter.ChapterNumber),
		}
	}

	return chapters
}

func MapToChapterDto(dbNovel *repositories.Novel, dbChapter *repositories.Chapter, prevChapter, nextChapter *int) chaptersdtostructs.ChapterPage {
	return chaptersdtostructs.ChapterPage{
		NovelName:      dbNovel.Title,
		NovelSlug:      dbNovel.Slug,
		ChapterTitle:   dbChapter.Title,
		ChapterNumber:  int(dbChapter.ChapterNumber),
		ChapterContent: dbChapter.Content,
		ChapterID:      int(dbChapter.ID),
		PrevChapter:    prevChapter,
		NextChapter:    nextChapter,
	}
}
