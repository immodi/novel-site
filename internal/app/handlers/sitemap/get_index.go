package sitemap

import (
	"net/http"
	"time"
)

func (s *SitemapHandler) HomeSiteMap(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/xml")

	// Get today's date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")

	// Home page sitemap XML
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
   <url>
      <loc>https://inovelhub.com/</loc>
      <lastmod>` + today + `</lastmod>
      <changefreq>daily</changefreq>
      <priority>1.0</priority>
   </url>
</urlset>`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml))
}
