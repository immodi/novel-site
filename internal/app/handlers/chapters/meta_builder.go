package chapters

import (
	"fmt"
	"immodi/novel-site/internal/db/repositories"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	"immodi/novel-site/pkg"
)

func BuildChapterMeta(dbNovel repositories.Novel, chapterNum int, novelStatus string) *indexdtostructs.MetaDataStruct {
	return &indexdtostructs.MetaDataStruct{
		IsRendering: true,
		Title:       fmt.Sprintf("%s - Chapter %d - Read %s Online - %s", dbNovel.Title, chapterNum, dbNovel.Title, indexdtostructs.SITE_NAME),
		Description: fmt.Sprintf("Read %s Chapter %d online. %s - Chapter %d by %s for free in high quality.", dbNovel.Title, chapterNum, dbNovel.Title, chapterNum, dbNovel.Author),
		Keywords:    fmt.Sprintf("read %s chapter %d online, free %s chapter %d, %s novel chapter %d", dbNovel.Title, chapterNum, dbNovel.Title, chapterNum, dbNovel.Title, chapterNum),
		OgURL:       fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, pkg.TitleToSlug(dbNovel.Title), chapterNum),
		Canonical:   fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, pkg.TitleToSlug(dbNovel.Title), chapterNum),

		// Extra
		CoverImage: fmt.Sprintf("%s/media/novel/%s.jpg", indexdtostructs.DOMAIN, pkg.TitleToSlug(dbNovel.Title)),
		Author:     dbNovel.Author,
		Status:     novelStatus,

		AuthorLink: fmt.Sprintf("%s/author/%s", indexdtostructs.DOMAIN, pkg.TitleToSlug(dbNovel.Author)),
		NovelName:  dbNovel.Title,

		// Navigation
		ReadURL:    fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, pkg.TitleToSlug(dbNovel.Title), chapterNum),
		UpdateTime: dbNovel.UpdateTime,
	}
}
