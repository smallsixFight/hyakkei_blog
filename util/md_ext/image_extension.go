package md_ext

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var cache sync.Map

func getNetImgWidthHeight(url string) (width, height int, err error) {
	if v, ok := cache.Load(url); ok {
		list := v.([]int64)
		if list != nil && len(list) == 3 {
			if time.Now().Unix()-list[2] < 3600 {
				return int(list[0]), int(list[1]), nil
			}
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("resp code: %d", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	if img, err := png.DecodeConfig(bytes.NewReader(data)); err == nil {
		height = img.Height
		width = img.Width
	} else if img, err := jpeg.DecodeConfig(bytes.NewReader(data)); err == nil {
		height = img.Height
		width = img.Width
	} else if img, err := gif.DecodeConfig(bytes.NewReader(data)); err == nil {
		height = img.Height
		width = img.Width
	} else {
		return 0, 0, fmt.Errorf("can not decode, please use jpeg/png/gif. ")
	}
	cache.Store(url, []int64{int64(width), int64(height), time.Now().Unix()})
	return
}

type imageExtension struct{}

var ImageExt = &imageExtension{}

func (p *imageExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(
			parser.NewLinkParser(), 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&myImageRenderer{}, 800),
		),
	)
}

type myImageRenderer struct {
	html.Renderer
}

func (r *myImageRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.render)
}

func (r *myImageRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString("<figure class=\"pswp-item\"><img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		_, _ = w.WriteString(`" `)
		if width, height, err := getNetImgWidthHeight(string(util.EscapeHTML(util.URLEscape(n.Destination, true)))); err == nil {
			_, _ = w.WriteString(fmt.Sprintf(`loading="lazy" width="%d" height="%d" `, width, height))
		}
	}
	_, _ = w.WriteString(`alt="`)
	_, _ = w.Write(util.EscapeHTML(n.Text(source)))
	_ = w.WriteByte('"')
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	if n.Title != nil {
		_, _ = w.WriteString(fmt.Sprintf("<figcaption>%s</figcaption>", n.Title))
	}
	_, _ = w.WriteString("</figure>")
	return ast.WalkSkipChildren, nil
}
