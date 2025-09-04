package novels

import (
	"fmt"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/db/repositories"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	novelsdtostructs "immodi/novel-site/internal/http/structs/novels"
	"immodi/novel-site/pkg"
)

func CastDbChaptersToInfoChapters(dbChapters []repositories.Chapter) []novelsdtostructs.Chapter {
	var chapters []novelsdtostructs.Chapter
	for _, dbChapter := range dbChapters {
		chapters = append(chapters, novelsdtostructs.Chapter{
			Title:  dbChapter.Title,
			Number: int(dbChapter.ChapterNumber),
		})
	}
	return chapters
}

func IncrementNovelViews(service novels.NovelService, novelId int64) {
	err := service.IncrementNovelViewCount(novelId)
	if err != nil {
		fmt.Printf("Failed to increment novel views: %v\n", err)
	}
}

func MapDBNovelToNovel(
	dbNovel repositories.Novel,
	genres, tags []string,
	novelStatus string,
	totalChapters int,
	currentPage int,
	chapters []novelsdtostructs.Chapter,
	lastChapter *repositories.Chapter,
) *novelsdtostructs.Novel {
	return &novelsdtostructs.Novel{
		Name:                dbNovel.Title,
		Description:         dbNovel.Description,
		Author:              dbNovel.Author,
		Genres:              genres,
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
	}
}

func MapNovelToMetaData(
	novel novelsdtostructs.Novel,
	novelStatus string,
) *indexdtostructs.MetaDataStruct {
	return &indexdtostructs.MetaDataStruct{
		IsRendering:       true,
		Title:             fmt.Sprintf("%s - Read %s For Free - %s", novel.Name, novel.Name, indexdtostructs.SITE_NAME),
		Description:       novel.Description,
		Keywords:          fmt.Sprintf("%s novel 2025, read %s online 2025, free %s novel", novel.Name, novel.Name, novel.Name),
		OgURL:             fmt.Sprintf("%s/novel/%s", indexdtostructs.DOMAIN, novel.Name),
		Canonical:         fmt.Sprintf("%s/novel/%s", indexdtostructs.DOMAIN, novel.Name),
		CoverImage:        novel.CoverImage,
		Genres:            novel.Genres,
		Author:            novel.Author,
		Status:            novelStatus,
		AuthorLink:        fmt.Sprintf("%s/author/%s", indexdtostructs.DOMAIN, novel.Author),
		NovelName:         novel.Name,
		ReadURL:           fmt.Sprintf("%s/novel/%s/chapter-1", indexdtostructs.DOMAIN, novel.Name),
		UpdateTime:        novel.LastUpdated,
		LatestChapterName: novel.LastChapterName,
		LatestChapterURL:  fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, novel.Name, novel.TotalChaptersNumber),
	}
}
