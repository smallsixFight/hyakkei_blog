package md_ext

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var StrikethroughExt = &strikethroughExtension{}

type strikethroughExtension struct{}

func (p *strikethroughExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(
			&strikethroughParser{}, 890),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&strikethroughHTMLRenderer{}, 900),
		),
	)
}

type strikethrough struct {
	ast.String
	val []byte
}

var kindStrikethrough = ast.NewNodeKind("delLine")

func (n *strikethrough) Kind() ast.NodeKind {
	return kindStrikethrough
}

type strikethroughParser struct{}

func (p *strikethroughParser) Trigger() []byte {
	return []byte("~~")
}

func (p *strikethroughParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 3 {
		return nil
	}
	pos := 2
	success := false
	for pos < len(line)-1 {
		if line[pos] == '~' && line[pos+1] == '~' {
			success = true
			break
		}
		pos++
	}
	if !success {
		return nil
	}
	block.Advance(pos + 2)
	return &strikethrough{val: line[2:pos]}
}

//  strikethroughHTMLRenderer 自定义渲染器
type strikethroughHTMLRenderer struct{}

func (r *strikethroughHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindStrikethrough, r.render)
}

func (r *strikethroughHTMLRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*strikethrough)
	if entering {
		_, _ = w.WriteString("<del>")
		_, _ = w.Write(n.val)
	} else {
		_, _ = w.WriteString("</del>")
	}
	return ast.WalkContinue, nil
}
