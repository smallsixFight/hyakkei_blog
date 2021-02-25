package md_ext

import (
	"bytes"
	"github.com/yuin/goldmark"
	"testing"
)

func TestDelLineExtension_Extend(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			StrikethroughExt,
		),
	)
	source := []byte("测试~~delete line~~ ~~aa~~")
	var b bytes.Buffer
	err := markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
	source = []byte("我来组成头部~~delete line~~我来组成躯干~~aa~~我来组成腿部")
	b.Reset()
	err = markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
}
