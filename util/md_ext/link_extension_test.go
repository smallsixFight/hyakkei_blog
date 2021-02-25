package md_ext

import (
	"bytes"
	"github.com/yuin/goldmark"
	"testing"
)

func TestLinkExtension_Extend(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			LinkExt,
		),
	)
	var b bytes.Buffer
	var source []byte
	source = []byte("[blog](https://blog.lamlake.com)")
	if err := markdown.Convert(source, &b); err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
	source = []byte("[github ](https://github.com/smallsixFight#{target=_blank})")
	b.Reset()
	if err := markdown.Convert(source, &b); err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
	source = []byte("[github ](https://github.com/smallsixFight#{target=_blank})后面跟着文字")
	b.Reset()
	if err := markdown.Convert(source, &b); err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
	source = []byte("[github ](https://github.com/smallsixFight#{target=_blank&title=title})多个属性")
	b.Reset()
	if err := markdown.Convert(source, &b); err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
}
