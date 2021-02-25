package md_ext

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"strings"
)

type linkExtension struct{}

var LinkExt = &linkExtension{}

func (p *linkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(
			&linkParser{}, 199), // 默认优先级为 200，所以自定义的需要高点
		),
	)
}

type linkParser struct{}

func (s *linkParser) Trigger() []byte {
	return []byte{'!', '[', ']'}
}

func (s *linkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	// 使用原来的解析器
	node := parser.NewLinkParser().Parse(parent, block, pc)
	n, ok := node.(*ast.Link)
	// 节点存在且为内容节点（即 node 非用于判断类型，node，'!'、'['会先解析为 node，再解析后面的内容）
	// 起始 node，不需进行 attr 处理。
	if node != nil && ok && len(line) > 0 {
		idx := strings.Index(string(n.Destination), "#{")
		if idx < 0 {
			return node
		}
		lr := text.NewReader(line)
		lr.Advance(len(n.Title) + 3 + idx)
		_, start := lr.Position()
		lr = text.NewReader(bytes.ReplaceAll(lr.Source()[start.Start:], []byte{'&'}, []byte{' '}))
		var attrs parser.Attributes
		var isOpen bool
		for {
			c := lr.Peek()
			if c == text.EOF || (c == '}' && isOpen) {
				break
			}
			if c == '\\' {
				lr.Advance(1)
				c = lr.Peek()
				if c == '{' || c == '}' {
					lr.Advance(1)
				}
				continue
			}
			if c == '{' {
				isOpen = true
				attrs, ok = parser.ParseAttributes(lr)
				if ok {
					lr.AdvanceLine()
					n.Destination = n.Destination[:idx]
				}
			}
			lr.Advance(1)
		}
		for i := range attrs {
			node.SetAttribute(attrs[i].Name, attrs[i].Value)
		}
	}
	return node
}
