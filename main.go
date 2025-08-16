package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/axadrn/axeladrian/assets"
	"github.com/axadrn/axeladrian/ui/pages"
	"github.com/joho/godotenv"
)

func main() {
	InitDotEnv()
	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)
	mux.Handle("GET /", templ.Handler(pages.Landing()))
	fmt.Println("Server is running on http://localhost:8090")

	http.ListenAndServe(":8090", mux)
}

func InitDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	isDev := os.Getenv("GO_ENV") != "production"

	// Main assets route
	var fs http.Handler
	if isDev {
		fs = http.FileServer(http.Dir("./assets"))
	} else {
		fs = http.FileServer(http.FS(assets.Assets))
	}
	
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDev {
			w.Header().Set("Cache-Control", "no-store")
		}
		fs.ServeHTTP(w, r)
	})))

	// Favicon routes for Safari compatibility
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
		// Set content type
		if strings.HasSuffix(path, ".ico") {
			w.Header().Set("Content-Type", "image/x-icon")
		} else {
			w.Header().Set("Content-Type", "image/png")
		}
		
		if isDev {
			w.Header().Set("Cache-Control", "no-store")
			http.ServeFile(w, r, "./assets/"+path)
		} else {
			content, err := assets.Assets.ReadFile(path)
			if err != nil {
				http.Error(w, "Favicon not found", http.StatusNotFound)
				return
			}
			w.Write(content)
		}
	}
}
