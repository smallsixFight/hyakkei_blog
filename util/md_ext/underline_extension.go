package md_ext

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var UnderlineExt = &underlineExtension{}

type underlineExtension struct{}

func (p *underlineExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(
			&underlineParser{}, 890),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&underlineHTMLRenderer{}, 900),
		),
	)
}

type underlineNode struct {
	ast.String
	val []byte
}

var kindUnderline = ast.NewNodeKind("underline")

func (n *underlineNode) Kind() ast.NodeKind {
	return kindUnderline
}

type underlineParser struct{}

func (p *underlineParser) Trigger() []byte {
	return []byte("++")
}

func (p *underlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 3 {
		return nil
	}
	pos := 2
	success := false
	for pos < len(line)-1 {
		if line[pos] == '+' && line[pos+1] == '+' {
			success = true
			break
		}
		pos++
	}
	if !success {
		return nil
	}
	block.Advance(pos + 2)
	return &underlineNode{val: line[2:pos]}
}

// 渲染生成 html 节点
type underlineHTMLRenderer struct{}

func (r *underlineHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindUnderline, r.render)
}

func (r *underlineHTMLRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*underlineNode)
	if entering {
		_, _ = w.WriteString("<ins>")
		_, _ = w.Write(n.val)
	} else {
		_, _ = w.WriteString("</ins>")
	}
	return ast.WalkContinue, nil
}
