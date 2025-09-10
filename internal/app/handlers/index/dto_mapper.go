package index

import (
	"immodi/novel-site/internal/app/services/index"
	"immodi/novel-site/internal/db/repositories"
	homenovelsdto "immodi/novel-site/internal/http/structs/index"
)

func DbNovelToHomeNovelMapper(dbNovels []repositories.Novel, homeService index.HomeService) ([]homenovelsdto.HomeNovelDto, error) {
	novels := make([]homenovelsdto.HomeNovelDto, 0, len(dbNovels))

	for _, dbNovel := range dbNovels {

		var dbLatestChapter repositories.Chapter

		dbLatestChapter, err := homeService.GetLatestChapterByNovel(dbNovel.ID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				dbLatestChapter = repositories.Chapter{
					Title: "Chapter doesn't exist",
				}
			} else {
				return nil, err
			}
		}

		dbNovelGenres, err := homeService.ListGenresByNovel(dbNovel.ID)
		if err != nil {
			return nil, err
		}

		dbChaptersCount, err := homeService.CountChaptersByNovel(dbNovel.ID)
		if err != nil {
			return nil, err
		}

		novelStatus := "Completed"
		if dbNovel.IsCompleted == 0 {
			novelStatus = "Ongoing"
		}

		novels = append(novels, MapToHomeNovelDto(dbNovel, dbLatestChapter, dbNovelGenres, dbChaptersCount, novelStatus))
	}

	return novels, nil
}

func MapToHomeNovelDto(
	dbNovel repositories.Novel,
	dbLatestChapter repositories.Chapter,
	dbNovelGenres []repositories.NovelGenre,
	dbChaptersCount int64,
	novelStatus string,
) homenovelsdto.HomeNovelDto {
	return homenovelsdto.HomeNovelDto{
		Name:                 dbNovel.Title,
		Slug:                 dbNovel.Slug,
		CoverImage:           dbNovel.CoverImage,
		LastestChapterNumber: int(dbLatestChapter.ChapterNumber),
		LastestChapterName:   dbLatestChapter.Title,
		Status:               novelStatus,
		Genres:               dbNovelGenres,
		LastUpdated:          dbNovel.UpdateTime,
		ChaptersCount:        int(dbChaptersCount),
	}
}
