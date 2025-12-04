package handler

import (
	"net/http"

	"github.com/axadrn/axeladrian/internal/service"
)

type SEOHandler struct {
	sitemapService *service.SitemapService
}

func NewSEOHandler(blogService *service.BlogService, baseURL string) *SEOHandler {
	return &SEOHandler{
		sitemapService: service.NewSitemapService(blogService, baseURL),
	}
}

func (h *SEOHandler) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(`User-agent: *
Allow: /
Sitemap: https://axeladrian.com/sitemap.xml`))
}

func (h *SEOHandler) Sitemap(w http.ResponseWriter, r *http.Request) {
	sitemap, err := h.sitemapService.GenerateSitemap()
	if err != nil {
		http.Error(w, "Failed to generate sitemap", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write(sitemap)
}
