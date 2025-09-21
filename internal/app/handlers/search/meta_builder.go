package search

import (
	"fmt"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	"net/url"
	"strings"
)

func BuildSearchMeta(
	searchCategory string,
	query string,
) *indexdtostructs.MetaDataStruct {
	searchCategory = strings.ToLower(searchCategory)
	query = strings.ToLower(query)
	searchURL := fmt.Sprintf("%s/%s/%s", indexdtostructs.DOMAIN, searchCategory, url.PathEscape(query))

	meta := &indexdtostructs.MetaDataStruct{
		IsRendering: true,
		Title:       fmt.Sprintf("Search results for \"%s\" in %s - %s", query, searchCategory, indexdtostructs.SITE_NAME),
		Description: fmt.Sprintf("Find free %s novels related to \"%s\". Browse thousands of light novels, web novels, and translated stories updated daily.", searchCategory, query),
		Keywords:    fmt.Sprintf("%s novels 2025, read %s novels online, free %s stories", query, query, query),
		OgURL:       searchURL,
		Canonical:   searchURL,
		Author:      indexdtostructs.SITE_NAME,
	}

	// For search/tag pages, structured data should describe the "search action"
	jsonLd := buildSearchJSONLD(query, searchURL, searchCategory)
	meta.NovelPageJsonLd = jsonLd

	return meta
}

func buildSearchJSONLD(query string, url string, searchCategory string,
) string {
	return fmt.Sprintf(`<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "SearchResultsPage",
  "name": "Search results for %s",
  "url": "%s",
  "about": "%s",
  "potentialAction": {
    "@type": "SearchAction",
    "target": "https://inovelhub.com/%s/{search_term_string}",
    "query-input": "required name=search_term_string"
  }
}
</script>`, query, url, query, searchCategory)
}
