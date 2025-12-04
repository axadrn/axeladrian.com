package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/axadrn/axeladrian/assets"
	"github.com/axadrn/axeladrian/content"
	"github.com/axadrn/axeladrian/internal/config"
	"github.com/axadrn/axeladrian/internal/handler"
	"github.com/axadrn/axeladrian/internal/middleware"
	"github.com/axadrn/axeladrian/internal/service"
	"github.com/axadrn/axeladrian/ui/pages"
	"github.com/joho/godotenv"
)

func main() {
	InitDotEnv()
	cfg := config.Load()

	// Content FS
	var contentFS fs.FS
	if cfg.IsDev() {
		contentFS = os.DirFS("content")
	} else {
		contentFS = content.Content
	}

	// Services
	blogService := service.NewBlogService(contentFS)

	// Handlers
	blogHandler := handler.NewBlogHandler(blogService)
	newsletterHandler := handler.NewNewsletterHandler(cfg)
	seoHandler := handler.NewSEOHandler(blogService, cfg.AppURL)

	mux := http.NewServeMux()

	// Assets
	SetupAssetsRoutes(mux, cfg)

	// Pages
	mux.Handle("GET /", templ.Handler(pages.Landing()))
	mux.HandleFunc("GET /blog", blogHandler.ListPosts)
	mux.HandleFunc("GET /blog/{slug}", blogHandler.ShowPost)
	mux.HandleFunc("GET /blog/tag/{tag}", blogHandler.ListByTag)

	// API
	mux.HandleFunc("POST /api/subscribe", newsletterHandler.Subscribe)

	// SEO
	mux.HandleFunc("GET /robots.txt", seoHandler.Robots)
	mux.HandleFunc("GET /sitemap.xml", seoHandler.Sitemap)

	// Apply middleware
	handler := middleware.WithConfig(cfg)(mux)

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", handler)
}

func InitDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func SetupAssetsRoutes(mux *http.ServeMux, cfg *config.Config) {
	isDev := cfg.IsDev()

	var fsHandler http.Handler
	if isDev {
		fsHandler = http.FileServer(http.Dir("./assets"))
	} else {
		fsHandler = http.FileServer(http.FS(assets.Assets))
	}

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDev {
			w.Header().Set("Cache-Control", "no-store")
		}
		fsHandler.ServeHTTP(w, r)
	})))

	favicons := map[string]string{
		"/favicon.ico":          "img/favicon/favicon.ico",
		"/apple-touch-icon.png": "img/favicon/apple-touch-icon.png",
		"/favicon-32x32.png":    "img/favicon/favicon-32x32.png",
		"/favicon-16x16.png":    "img/favicon/favicon-16x16.png",
	}

	for route, path := range favicons {
		mux.HandleFunc("GET "+route, serveFavicon(path, isDev))
	}
}

func serveFavicon(path string, isDev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(path, ".ico") {
			w.Header().Set("Content-Type", "image/x-icon")
		} else {
			w.Header().Set("Content-Type", "image/png")
		}

		if isDev {
			w.Header().Set("Cache-Control", "no-store")
			http.ServeFile(w, r, "./assets/"+path)
		} else {
			fileContent, err := assets.Assets.ReadFile(path)
			if err != nil {
				http.Error(w, "Favicon not found", http.StatusNotFound)
				return
			}
			w.Write(fileContent)
		}
	}
}
