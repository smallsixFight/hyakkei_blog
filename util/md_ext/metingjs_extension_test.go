package md_ext

import (
	"bytes"
	"github.com/yuin/goldmark"
	"testing"
)

func TestMetingExtension_Extend(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			MetingJS,
		),
	)
	source := []byte("#meting server=srv id=id typ=song")
	var b bytes.Buffer
	err := markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
	source = []byte("#meting server=srv id=id typ=song\n\n#meting server=\"srv\" id=\"id\" typ=\"song\"")
	b.Reset()
	err = markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
}
