package sitemap

import (
	"fmt"
	"immodi/novel-site/pkg"
	"net/http"
	"time"
)

func (s *SitemapHandler) NovelsSiteMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	dbNovels, err := s.sitemapService.GetAllNovels()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// build the <urlset>
	xml := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"

	for _, n := range dbNovels {
		lastmod, err := pkg.GetDateFromRFCStrDash(n.UpdateTime)
		if err != nil {
			lastmod = time.Now().UTC().Format("2006-01-02")
		}

		xml += fmt.Sprintf(
			"  <url>\n"+
				"    <loc>https://inovelhub.com/novel/%s</loc>\n"+
				"    <lastmod>%s</lastmod>\n"+
				"  </url>\n",
			n.Slug, lastmod)
	}

	xml += `</urlset>`

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml))
}
