package filter

import (
	"fmt"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
)

func BuildFilterMeta() *indexdtostructs.MetaDataStruct {
	return &indexdtostructs.MetaDataStruct{
		IsRendering:    true,
		Title:          fmt.Sprintf("Browse Novels by Tags & Genres - %s", indexdtostructs.SITE_NAME),
		Description:    "Discover thousands of free light novels organized by tags and genres. Explore fantasy, romance, action, cultivation, and more—updated daily for your reading pleasure.",
		Keywords:       "light novel genres, webnovel tags, fantasy novels, romance novels, cultivation novels, free light novel",
		OgURL:          fmt.Sprintf("%s/filter", indexdtostructs.DOMAIN),
		Canonical:      fmt.Sprintf("%s/filter", indexdtostructs.DOMAIN),
		CoverImage:     fmt.Sprintf("%s/static/logo/logo.png", indexdtostructs.DOMAIN),
		Author:         indexdtostructs.SITE_NAME,
		HomePageJsonLd: filterJsonLd(),
	}
}

func filterJsonLd() string {
	return `<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "CollectionPage",
  "name": "INovelHub - Browse Novels by Tags & Genres",
  "url": "https://inovelhub.com/filter",
  "description": "Browse and read free light novels filtered by tags and genres such as fantasy, romance, action and more.",
  "isPartOf": {
    "@type": "WebSite",
    "name": "INovelHub",
    "url": "https://inovelhub.com"
  }
}
</script>`
}
