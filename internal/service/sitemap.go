package service

import (
	"encoding/xml"
	"log/slog"
	"strings"
	"time"

	"github.com/axadrn/axeladrian/internal/model"
)

var publicRoutes = []struct {
	Path       string
	Priority   string
	ChangeFreq string
}{
	{"/", "1.0", "weekly"},
	{"/blog", "0.8", "daily"},
}

type SitemapService struct {
	blogService *BlogService
	baseURL     string
}

func NewSitemapService(blogService *BlogService, baseURL string) *SitemapService {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &SitemapService{
		blogService: blogService,
		baseURL:     baseURL,
	}
}

func (s *SitemapService) GenerateSitemap() ([]byte, error) {
	sitemap := model.Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []model.SitemapURL{},
	}

	// Add static routes
	for _, route := range publicRoutes {
		sitemap.URLs = append(sitemap.URLs, model.SitemapURL{
			Loc:        s.baseURL + route.Path,
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: route.ChangeFreq,
			Priority:   route.Priority,
		})
	}

	// Add blog posts
	posts, err := s.blogService.Posts()
	if err != nil {
		slog.Warn("failed to get blog URLs for sitemap", "error", err)
	} else {
		for _, post := range posts {
			lastMod := time.Now().Format("2006-01-02")
			if !post.Date.IsZero() {
				lastMod = post.Date.Format("2006-01-02")
			}
			sitemap.URLs = append(sitemap.URLs, model.SitemapURL{
				Loc:        s.baseURL + "/blog/" + post.Slug,
				LastMod:    lastMod,
				ChangeFreq: "weekly",
				Priority:   "0.7",
			})
		}

		// Add tag pages
		tagMap := make(map[string]bool)
		for _, post := range posts {
			for _, tag := range post.Tags {
				tagMap[tag] = true
			}
		}
		for tag := range tagMap {
			sitemap.URLs = append(sitemap.URLs, model.SitemapURL{
				Loc:        s.baseURL + "/blog/tag/" + tag,
				LastMod:    time.Now().Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.5",
			})
		}
	}

	output, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return nil, err
	}

	return []byte(xml.Header + string(output)), nil
}
