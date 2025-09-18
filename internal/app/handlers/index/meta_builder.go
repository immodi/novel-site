package index

import (
	"fmt"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
)

func BuildHomeMeta() *indexdtostructs.MetaDataStruct {
	return &indexdtostructs.MetaDataStruct{
		IsRendering:    true,
		Title:          fmt.Sprintf("Read Free Light Novel Online - %s", indexdtostructs.SITE_NAME),
		Description:    "We are offering thousands of free books online read! Read novel updated daily: light novel translations, web novel, chinese novel, japanese novel, korean novel, english novel and other novels online.",
		Keywords:       "freewebnovel, novellive, novelfull, mtlnovel, novelupdates, webnovel, korean novel, cultivation novel",
		OgURL:          indexdtostructs.DOMAIN,
		Canonical:      indexdtostructs.DOMAIN,
		CoverImage:     fmt.Sprintf("%s/static/logo/logo.png", indexdtostructs.DOMAIN),
		Author:         indexdtostructs.SITE_NAME,
		HomePageJsonLd: homeJsonLd(),
	}
}

func homeJsonLd() string {
	return `<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "INovelHub",
  "url": "https://inovelhub.com",
  "potentialAction": {
    "@type": "SearchAction",
    "target": "https://inovelhub.com/search/{search_term_string}",
    "query-input": "required name=search_term_string"
  }
}
</script>`
}
