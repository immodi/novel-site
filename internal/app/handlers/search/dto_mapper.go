package search

import (
	"immodi/novel-site/internal/app/services/search"
	"immodi/novel-site/internal/db/repositories"
	searchresutlsdto "immodi/novel-site/internal/http/structs/search"
	"immodi/novel-site/pkg"
)

func DbNovelsToSearchDtosMapper(service search.SearchService, dbNovels []repositories.Novel) []searchresutlsdto.SearchResultDto {
	var searchDtos []searchresutlsdto.SearchResultDto
	for _, novel := range dbNovels {
		status := "Ongoing"
		if novel.IsCompleted == 1 {
			status = "Completed"
		}

		lastChapterName := "Chapter not found"
		lastChapter, err := service.GetLastChapter(novel.ID)
		if err == nil {
			lastChapterName = "Couldn't find last chapter name"
		}

		if lastChapter.ID != 0 {
			lastChapterName = lastChapter.Title
		}

		searchDtos = append(searchDtos, searchresutlsdto.SearchResultDto{
			ID:                   int(novel.ID),
			Name:                 novel.Title,
			Slug:                 novel.Slug,
			CoverImage:           novel.CoverImage,
			Author:               novel.Author,
			Status:               status,
			LastestChapterName:   lastChapterName,
			LastestChapterNumber: int(lastChapter.ChapterNumber),
			LastUpdated:          pkg.TimeAgo(novel.UpdateTime),
		})
	}
	return searchDtos

}
