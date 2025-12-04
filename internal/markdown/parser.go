package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

type Parser struct {
	md goldmark.Markdown
}

func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
		),
	)

	return &Parser{md: md}
}

func (p *Parser) Parse(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := p.md.Convert(source, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Parser) ParseWithFrontmatter(source []byte) (content []byte, meta map[string]any, err error) {
	context := parser.NewContext()
	var buf bytes.Buffer

	err = p.md.Convert(source, &buf, parser.WithContext(context))
	if err != nil {
		return nil, nil, err
	}

	data := frontmatter.Get(context)
	if data == nil {
		meta = make(map[string]any)
	} else {
		err = data.Decode(&meta)
		if err != nil {
			meta = make(map[string]any)
		}
	}

	return buf.Bytes(), meta, nil
}
