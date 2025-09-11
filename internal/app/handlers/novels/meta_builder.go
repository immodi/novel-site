package novels

import (
	"fmt"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	novelsdtostructs "immodi/novel-site/internal/http/structs/novels"
)

func BuildNovelMeta(
	novel *novelsdtostructs.Novel,
	novelStatus string,
) *indexdtostructs.MetaDataStruct {
	var genres []string
	for _, g := range novel.Genres {
		genres = append(genres, g.Genre)
	}
	return &indexdtostructs.MetaDataStruct{
		IsRendering:       true,
		Title:             fmt.Sprintf("%s - Read %s For Free - %s", novel.Name, novel.Name, indexdtostructs.SITE_NAME),
		Description:       novel.Description,
		Keywords:          fmt.Sprintf("%s novel 2025, read %s online 2025, free %s novel", novel.Name, novel.Name, novel.Name),
		OgURL:             fmt.Sprintf("%s/novel/%s", indexdtostructs.DOMAIN, novel.Name),
		Canonical:         fmt.Sprintf("%s/novel/%s", indexdtostructs.DOMAIN, novel.Name),
		CoverImage:        novel.CoverImage,
		Genres:            genres,
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
