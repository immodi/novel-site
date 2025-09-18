package chapters

import (
	"fmt"
	"immodi/novel-site/internal/db/repositories"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	"immodi/novel-site/pkg"
	"strings"
)

func BuildChapterMeta(dbNovel *repositories.Novel, chapterNum int, novelStatus string) *indexdtostructs.MetaDataStruct {
	meta := &indexdtostructs.MetaDataStruct{
		IsRendering: true,
		Title:       fmt.Sprintf("%s - Chapter %d - Read %s Online - %s", dbNovel.Title, chapterNum, dbNovel.Title, indexdtostructs.SITE_NAME),
		Description: fmt.Sprintf("Read %s Chapter %d online. %s - Chapter %d by %s for free in high quality.", dbNovel.Title, chapterNum, dbNovel.Title, chapterNum, dbNovel.Author),
		Keywords:    fmt.Sprintf("read %s chapter %d online, free %s chapter %d, %s novel chapter %d", dbNovel.Title, chapterNum, dbNovel.Title, chapterNum, dbNovel.Title, chapterNum),
		OgURL:       fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, dbNovel.Slug, chapterNum),
		Canonical:   fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, dbNovel.Slug, chapterNum),

		// Extra
		CoverImage: fmt.Sprintf("%s%s", indexdtostructs.DOMAIN, dbNovel.CoverImage),
		Author:     dbNovel.Author,
		Status:     novelStatus,

		AuthorLink: fmt.Sprintf("%s/author/%s", indexdtostructs.DOMAIN, dbNovel.AuthorSlug),
		NovelName:  dbNovel.Title,

		// Navigation
		ReadURL:    fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, dbNovel.Slug, chapterNum),
		UpdateTime: dbNovel.UpdateTime,
	}

	meta.ChapterPageJsonLd = buildChapterJSONLD(meta, chapterNum)
	return meta
}

func buildChapterJSONLD(meta *indexdtostructs.MetaDataStruct, chapterNum int) string {
	// Convert UpdateTime to dash-format YYYY-MM-DD
	date, err := pkg.GetDateFromRFCStrDash(meta.UpdateTime)
	if err != nil {
		date = meta.UpdateTime
	}

	jsonld := fmt.Sprintf(`{
  "@context": "https://schema.org",
  "@type": "Chapter",
  "name": "Chapter %d",
  "url": "%s",
  "partOf": {
    "@type": "Book",
    "name": "%s",
    "author": {
      "@type": "Person",
      "name": "%s",
      "url": "%s"
    },
    "publisher": {
      "@type": "Organization",
      "name": "%s"
    },
    "url": "%s",
    "image": "%s",
    "genre": "%s",
    "datePublished": "%s",
    "inLanguage": "en"
  },
  "datePublished": "%s"
}`,
		chapterNum,
		meta.OgURL,     // chapter URL
		meta.NovelName, // Book name
		meta.Author,
		meta.AuthorLink,
		indexdtostructs.SITE_NAME,
		meta.OgURL, // Book URL
		meta.CoverImage,
		strings.Join(meta.Genres, ", "),
		date, // Book published date
		date, // Chapter published date
	)

	return fmt.Sprintf(`<script type="application/ld+json">%s</script>`, jsonld)
}
