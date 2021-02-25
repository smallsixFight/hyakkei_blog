package md_ext

import (
	"bytes"
	"github.com/yuin/goldmark"
	"testing"
)

func TestImageExtension_Extend(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			ImageExt,
		),
	)
	source := []byte("![img_1](https://i.loli.net/2021/01/20/uWnCAXSbRO9Q6Zr.png)")
	var b bytes.Buffer
	err := markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
	source = []byte("![img_2](https://i.loli.net/2021/01/20/uWnCAXSbRO9Q6Zr.png)")
	b.Reset()
	err = markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
}
