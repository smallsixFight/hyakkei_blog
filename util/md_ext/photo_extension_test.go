package md_ext

import (
	"bytes"
	"github.com/yuin/goldmark"
	"image/png"
	"net/http"
	"testing"
)

func TestGetPhotoHW(t *testing.T) {
	resp, err := http.Get("https://i.loli.net/2021/01/20/uWnCAXSbRO9Q6Zr.png")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	cfg, err := png.DecodeConfig(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg.Width, cfg.Height)
}

func TestPhotoExtension_Extend(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Photos,
			ImageExt,
		),
	)
	source := []byte("[photos]\n![test](https://static.imalan.cn/file/image/post/79ade.png \"Optional title\")\n![test](https://i.loli.net/2021/01/20/lUFXK4xO6JWzgS3.jpg)\n[/photos]")
	var b bytes.Buffer
	err := markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b.String())
}
