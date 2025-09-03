package index

import (
	"fmt"
	"immodi/novel-site/internal/db/repositories"
	homenovelsdto "immodi/novel-site/internal/http/structs/index"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
)

func dbNovelToHomeNovelMapper(dbNovels []repositories.Novel, h *HomeHandler) ([]homenovelsdto.HomeNovelDto, error) {
	novels := make([]homenovelsdto.HomeNovelDto, 0, len(dbNovels))

	for _, dbNovel := range dbNovels {

		var dbLatestChapter repositories.Chapter

		dbLatestChapter, err := h.homeService.GetLatestChapterByNovel(dbNovel.ID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				dbLatestChapter = repositories.Chapter{
					Title: "Chapter doesn't exist",
				}
			} else {
				return nil, err
			}
		}

		dbNovelGenres, err := h.homeService.ListGenresByNovel(dbNovel.ID)
		if err != nil {
			return nil, err
		}

		dbChaptersCount, err := h.homeService.CountChaptersByNovel(dbNovel.ID)
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
	dbNovelGenres []string,
	dbChaptersCount int64,
	novelStatus string,
) homenovelsdto.HomeNovelDto {
	return homenovelsdto.HomeNovelDto{
		Name:                 dbNovel.Title,
		CoverImage:           dbNovel.CoverImage,
		LastestChapterNumber: int(dbLatestChapter.ChapterNumber),
		LastestChapterName:   dbLatestChapter.Title,
		Status:               novelStatus,
		Genres:               dbNovelGenres,
		LastUpdated:          dbNovel.UpdateTime,
		ChaptersCount:        int(dbChaptersCount),
	}
}

func GetIndexMetaData() *indexdtostructs.MetaDataStruct {
	return &indexdtostructs.MetaDataStruct{
		IsRendering: true,
		Title:       fmt.Sprintf("Read Free Light Novel Online - %s", indexdtostructs.SITE_NAME),
		Description: "We are offering thousands of free books online read! Read novel updated daily: light novel translations, web novel, chinese novel, japanese novel, korean novel, english novel and other novels online.",
		Keywords:    "freewebnovel, novellive, novelfull, mtlnovel, novelupdates, webnovel, korean novel, cultivation novel",
		OgURL:       indexdtostructs.DOMAIN,
		Canonical:   indexdtostructs.DOMAIN,
		CoverImage:  fmt.Sprintf("%s/img/cover.jpg", indexdtostructs.DOMAIN),
		Author:      indexdtostructs.SITE_NAME,
	}
}
