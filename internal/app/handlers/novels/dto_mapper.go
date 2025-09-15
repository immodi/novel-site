package novels

import (
	"immodi/novel-site/internal/db/repositories"
	novelsdtostructs "immodi/novel-site/internal/http/structs/novels"
	"immodi/novel-site/pkg"
)

func DbChaptersToInfoChaptersMapper(dbChapters []repositories.Chapter) []novelsdtostructs.Chapter {
	var chapters []novelsdtostructs.Chapter
	for _, dbChapter := range dbChapters {
		chapters = append(chapters, novelsdtostructs.Chapter{
			Title:       dbChapter.Title,
			ReleaseDate: pkg.TimeAgo(dbChapter.ReleaseDate),
			Number:      int(dbChapter.ChapterNumber),
		})
	}
	return chapters
}

func MapToNovel(
	dbNovel repositories.Novel,
	genres []repositories.NovelGenre,
	tags []repositories.NovelTag,
	novelStatus string,
	totalChapters int,
	isNovelBookMarked bool,
	currentPage int,
	chapters []novelsdtostructs.Chapter,
	lastChapter *repositories.Chapter,
	successMessage string,
	errorMessage string,
) *novelsdtostructs.Novel {
	return &novelsdtostructs.Novel{
		ID:                  int(dbNovel.ID),
		Name:                dbNovel.Title,
		Description:         dbNovel.Description,
		Author:              dbNovel.Author,
		Slug:                dbNovel.Slug,
		AuthorSlug:          dbNovel.AuthorSlug,
		Genres:              genres,
		Views:               pkg.AbbreviateInt(int(dbNovel.ViewCount)),
		Tags:                tags,
		Status:              novelStatus,
		ReleaseYear:         int(dbNovel.ReleaseYear),
		Publisher:           dbNovel.Publisher,
		CoverImage:          dbNovel.CoverImage,
		TotalChaptersNumber: totalChapters,
		CurrentPage:         currentPage,
		TotalPages:          pkg.CalculateTotalPages(totalChapters),
		Chapters:            chapters,
		LastChapterName:     lastChapter.Title,
		LastUpdated:         dbNovel.UpdateTime,
		IsNovelBookMarked:   isNovelBookMarked,

		SuccessMessage: successMessage,
		ErrorMessage:   errorMessage,
	}
}
