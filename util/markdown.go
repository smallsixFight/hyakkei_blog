package util

import (
	"github.com/smallsixFight/hyakkei_blog/util/md_ext"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"sync"
)

var mdHandle goldmark.Markdown
var once sync.Once

func GetMDHandle() goldmark.Markdown {
	once.Do(func() {
		mdHandle = goldmark.New(
			goldmark.WithExtensions(
				md_ext.Photos,
				md_ext.MetingJS,
				md_ext.UnderlineExt,
				md_ext.StrikethroughExt,
				md_ext.LinkExt,
				md_ext.ImageExt,
				extension.Table,
				extension.DefinitionList,
			),
		)
	})
	return mdHandle
}
