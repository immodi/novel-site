package sitemap

import (
	"fmt"
	"net/http"
	"time"
)

func (s *SitemapHandler) MainSiteMap(w http.ResponseWriter, r *http.Request) {
	// set correct content type
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	// dynamic lastmod date (UTC in YYYY-MM-DD)
	today := time.Now().UTC().Format("2006-01-02")

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>https://inovelhub.com/sitemaps/novels.xml</loc>
    <lastmod>%s</lastmod>
  </sitemap>
  <sitemap>
    <loc>https://inovelhub.com/sitemaps/genres.xml</loc>
    <lastmod>%s</lastmod>
  </sitemap>
  <sitemap>
    <loc>https://inovelhub.com/sitemaps/tags.xml</loc>
    <lastmod>%s</lastmod>
  </sitemap>
</sitemapindex>`, today, today, today)

	// write the XML
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml))
}
