package sitemap

import (
	"fmt"
	"net/http"
	"time"
)

func (s *SitemapHandler) TagsSiteMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	// fetch all tag slugs from the DB/service
	tags, err := s.sitemapService.GetAllTags()
	if err != nil {
		// log the real error on the server, return generic text to the client
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// build the <urlset>
	xml := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"

	// use current UTC date for lastmod
	lastmod := time.Now().UTC().Format("2006-01-02")

	for _, slug := range tags {
		xml += fmt.Sprintf(
			"  <url>\n"+
				"    <loc>https://inovelhub.com/tag/%s</loc>\n"+
				"    <lastmod>%s</lastmod>\n"+
				"  </url>\n",
			slug, lastmod)
	}

	xml += `</urlset>`

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml))
}
