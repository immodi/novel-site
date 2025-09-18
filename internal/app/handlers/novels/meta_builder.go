package novels

import (
	"fmt"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	novelsdtostructs "immodi/novel-site/internal/http/structs/novels"
	"immodi/novel-site/pkg"
	"strings"
)

func BuildNovelMeta(
	novel *novelsdtostructs.Novel,
	novelStatus string,
) *indexdtostructs.MetaDataStruct {
	var genres []string
	for _, g := range novel.Genres {
		genres = append(genres, g.Genre)
	}
	meta := &indexdtostructs.MetaDataStruct{
		IsRendering:       true,
		Title:             fmt.Sprintf("%s - Read %s For Free - %s", novel.Name, novel.Name, indexdtostructs.SITE_NAME),
		Description:       novel.Description,
		Keywords:          fmt.Sprintf("%s novel 2025, read %s online 2025, free %s novel", novel.Name, novel.Name, novel.Name),
		OgURL:             fmt.Sprintf("%s/novel/%s", indexdtostructs.DOMAIN, novel.Slug),
		Canonical:         fmt.Sprintf("%s/novel/%s", indexdtostructs.DOMAIN, novel.Slug),
		CoverImage:        fmt.Sprintf("%s%s", indexdtostructs.DOMAIN, novel.CoverImage),
		Genres:            genres,
		Author:            novel.Author,
		Status:            novelStatus,
		AuthorLink:        fmt.Sprintf("%s/author/%s", indexdtostructs.DOMAIN, novel.AuthorSlug),
		NovelName:         novel.Name,
		ReadURL:           fmt.Sprintf("%s/novel/%s/chapter-1", indexdtostructs.DOMAIN, novel.Slug),
		UpdateTime:        novel.LastUpdated,
		LatestChapterName: novel.LastChapterName,
		LatestChapterURL:  fmt.Sprintf("%s/novel/%s/chapter-%d", indexdtostructs.DOMAIN, novel.Slug, novel.TotalChaptersNumber),
	}
	jsonLd := buildNovelJSONLD(meta)
	meta.NovelPageJsonLd = jsonLd
	return meta
}

func buildNovelJSONLD(meta *indexdtostructs.MetaDataStruct) string {
	date, err := pkg.GetDateFromRFCStrDash(meta.UpdateTime)
	if err != nil {
		date = meta.UpdateTime
	}

	jsonld := fmt.Sprintf(`{
  "@context": "https://schema.org",
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
  "inLanguage": "en",
  "mainEntityOfPage": "%s"
}`,
		meta.NovelName,
		meta.Author,
		meta.AuthorLink,
		indexdtostructs.SITE_NAME,
		meta.OgURL,
		meta.CoverImage,
		strings.Join(meta.Genres, ", "),
		date,
		meta.OgURL,
	)

	// Wrap in <script> tag
	return fmt.Sprintf(`<script type="application/ld+json">%s</script>`, jsonld)
}
