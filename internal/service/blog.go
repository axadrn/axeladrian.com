package service

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"github.com/axadrn/axeladrian/internal/markdown"
	"github.com/axadrn/axeladrian/internal/model"
)

type BlogService struct {
	parser *markdown.Parser
	fs     fs.FS
}

func NewBlogService(contentFS fs.FS) *BlogService {
	return &BlogService{
		parser: markdown.NewParser(),
		fs:     contentFS,
	}
}

func (s *BlogService) Posts() ([]*model.BlogPost, error) {
	files, err := fs.Glob(s.fs, "blog/*.md")
	if err != nil {
		return nil, err
	}

	var posts []*model.BlogPost
	for _, file := range files {
		slug := strings.TrimSuffix(strings.TrimPrefix(file, "blog/"), ".md")
		post, err := s.Post(slug)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

func (s *BlogService) Post(slug string) (*model.BlogPost, error) {
	content, err := fs.ReadFile(s.fs, "blog/"+slug+".md")
	if err != nil {
		return nil, fmt.Errorf("blog post not found: %s", slug)
	}

	htmlContent, meta, err := s.parser.ParseWithFrontmatter(content)
	if err != nil {
		return nil, err
	}

	post := &model.BlogPost{
		Slug:        slug,
		HTMLContent: string(htmlContent),
		Content:     string(content),
	}

	if title, ok := meta["title"].(string); ok {
		post.Title = title
	}
	if author, ok := meta["author"].(string); ok {
		post.Author = author
	}
	if description, ok := meta["description"].(string); ok {
		post.Description = description
	}
	if image, ok := meta["image"].(string); ok {
		post.Image = image
	}
	if dateStr, ok := meta["date"].(string); ok {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			post.Date = date
		}
	}
	if tags, ok := meta["tags"].([]any); ok {
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				post.Tags = append(post.Tags, tagStr)
			}
		}
	}

	post.ReadTime = s.calculateReadTime(string(content))

	return post, nil
}

func (s *BlogService) PostsByTag(tag string) ([]*model.BlogPost, error) {
	allPosts, err := s.Posts()
	if err != nil {
		return nil, err
	}

	var posts []*model.BlogPost
	for _, post := range allPosts {
		for _, postTag := range post.Tags {
			if strings.EqualFold(postTag, tag) {
				posts = append(posts, post)
				break
			}
		}
	}

	return posts, nil
}

func (s *BlogService) calculateReadTime(content string) int {
	words := strings.Fields(content)
	wordsPerMinute := 200
	readTime := len(words) / wordsPerMinute
	if readTime < 1 {
		readTime = 1
	}
	return readTime
}
