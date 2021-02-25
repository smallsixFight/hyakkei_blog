package md_ext

import (
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"strings"
)

type metingJS struct{}

var MetingJS = &metingJS{}

func (p *metingJS) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(
			&metingJSParser{}, 700),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&metingJSRenderer{}, 700),
		),
	)
}

type metingJSParser struct{}

func (p *metingJSParser) Trigger() []byte {
	return []byte("#meting ")
}

func (p *metingJSParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	str := string(line)
	idx1, idx2, idx3 := strings.Index(str, " id="), strings.Index(str, " typ="), strings.Index(str, " server=")
	if idx1 < 0 || idx2 < 0 || idx3 < 0 {
		return nil
	}
	var id, typ, server strings.Builder
	for i := idx1 + 4; i < len(str); i++ {
		if str[i] == '"' {
			continue
		}
		if str[i] == ' ' || str[i] == '"' {
			break
		}
		id.WriteByte(str[i])
	}
	for i := idx2 + 5; i < len(str); i++ {
		if str[i] == '"' {
			continue
		}
		if str[i] == ' ' || str[i] == '"' {
			break
		}
		typ.WriteByte(str[i])
	}
	for i := idx3 + 8; i < len(str); i++ {
		if str[i] == '"' {
			continue
		}
		if str[i] == ' ' || str[i] == '"' {
			break
		}
		server.WriteByte(str[i])
	}
	block.AdvanceLine()
	return newMetingJSNode(id.String(), server.String(), typ.String())
}

func newMetingJSNode(id, server, typ string) *metingJSNode {
	c := &metingJSNode{
		id:     id,
		server: server,
		typ:    typ,
	}

	return c
}

type metingJSNode struct {
	ast.Link
	server string
	typ    string
	id     string
}

var kindMetingJS = ast.NewNodeKind("meting")

func (n *metingJSNode) Kind() ast.NodeKind {
	return kindMetingJS
}

type metingJSRenderer struct{}

func (r *metingJSRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindMetingJS, r.render)
}

func (r *metingJSRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*metingJSNode)
	if !entering {
		_, _ = w.WriteString(fmt.Sprintf(`<meting-js server="%s" type="%s" id="%s"></meting-js>`, n.server, n.typ, n.id))
	}
	return ast.WalkContinue, nil
}
