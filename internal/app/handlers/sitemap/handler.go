package sitemap

import "immodi/novel-site/internal/app/services/sitemap"

type SitemapHandler struct {
	sitemapService sitemap.SitemapService
}

func NewSitemapHandler(service sitemap.SitemapService) *SitemapHandler {
	return &SitemapHandler{sitemapService: service}
}
