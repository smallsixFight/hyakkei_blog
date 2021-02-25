package md_ext

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var Photos = &photos{}

type photos struct{}

func (p *photos) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(util.Prioritized(
			&photosParser{}, 890),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&photosHTMLRenderer{}, 990),
		),
	)
}

type photosNode struct {
	ast.HTMLBlock
}

var kindPhotos = ast.NewNodeKind("photos")

func (n *photosNode) Kind() ast.NodeKind {
	return kindPhotos
}

type photosParser struct{}

func (p *photosParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	// 头部匹配
	if string(line[:8]) != "[photos]" {
		return nil, parser.NoChildren
	}
	pos := 8
	reader.Advance(pos)
	return &photosNode{}, parser.HasChildren
}

func (p *photosParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	// 结尾匹配
	if len(line) > 8 && string(line[:9]) == "[/photos]" {
		reader.Advance(9)
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (p *photosParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

func (p *photosParser) CanInterruptParagraph() bool {
	return true
}

func (p *photosParser) CanAcceptIndentedLine() bool {
	return false
}

func (p *photosParser) Trigger() []byte {
	return []byte("[photos]")
}

// 渲染生成 html 节点
type photosHTMLRenderer struct{}

func (r *photosHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindPhotos, r.render)
}

func (r *photosHTMLRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<div class=\"photos\">\n")
	} else {
		_, _ = w.WriteString("</div>\n")
	}
	return ast.WalkContinue, nil
}
