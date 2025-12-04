package handler

import (
	"net/http"

	"github.com/axadrn/axeladrian/internal/service"
	"github.com/axadrn/axeladrian/internal/ui"
	"github.com/axadrn/axeladrian/ui/pages"
)

type BlogHandler struct {
	blogService *service.BlogService
}

func NewBlogHandler(blogService *service.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}

func (h *BlogHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.blogService.Posts()
	if err != nil {
		http.Error(w, "Failed to load blog posts", http.StatusInternalServerError)
		return
	}

	ui.Render(w, r, pages.BlogList(posts))
}

func (h *BlogHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	post, err := h.blogService.Post(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ui.Render(w, r, pages.BlogPost(post))
}

func (h *BlogHandler) ListByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	if tag == "" {
		http.NotFound(w, r)
		return
	}

	posts, err := h.blogService.PostsByTag(tag)
	if err != nil {
		http.Error(w, "Failed to load blog posts", http.StatusInternalServerError)
		return
	}

	ui.Render(w, r, pages.BlogListByTag(posts, tag))
}
