package file_generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var pageGenerateLock sync.Mutex

func GenerateBasePage(src, dest, pageName string) error {
	template, err := getTemplate(filepath.Join(src, pageName))
	if err != nil {
		return fmt.Errorf("读取 %s 模版文件数据失败: %s", pageName, err.Error())
	}
	str := ReplaceBasePath(template, "")
	return util.WriteFile([]byte(str), filepath.Join(dest, pageName))
}

func ReplaceBasePath(d []byte, preDir string) string {
	str := string(d)
	str = strings.ReplaceAll(str, "#{{img_path}}", filepath.Join(preDir, util.GetInitConfig().ImgPath))
	str = strings.ReplaceAll(str, "#{{js_path}}", filepath.Join(preDir, util.GetInitConfig().JsPath))
	str = strings.ReplaceAll(str, "#{{css_path}}", filepath.Join(preDir, util.GetInitConfig().CssPath))
	str = strings.ReplaceAll(str, "#{{font_path}}", filepath.Join(preDir, util.GetInitConfig().FontPath))
	str = strings.Replace(str, "#{{github_name}}", service.GetSysConfig().GithubName, 1)
	str = strings.Replace(str, "#{{blog_name}}", service.GetSysConfig().BlogName, 1)
	str = strings.Replace(str, "#{{pswp}}", pswpStr, 1)
	return str
}

// 生成初始页
func GenerateIndexPage() error {
	return GenerateBasePage(filepath.Join(util.GetBlogTemplatePath(), "blog_pages"), filepath.Join(service.GetSysConfig().SavePath, "hyakkei"), "index.html")
}

// 生成书籍页
func GenerateBookPage() error {
	return GenerateBasePage(filepath.Join(util.GetBlogTemplatePath(), "blog_pages"), filepath.Join(service.GetSysConfig().SavePath, "hyakkei"), "books.html")
}

// 生成友链页
func GenerateFriendLinksPage() error {
	return GenerateBasePage(filepath.Join(util.GetBlogTemplatePath(), "blog_pages"), filepath.Join(service.GetSysConfig().SavePath, "hyakkei"), "friends.html")
}

// 为文章生成对应 html 文件
func GenerateArticlePage(article *model.Post) error {
	pageGenerateLock.Lock()
	defer pageGenerateLock.Unlock()
	template, err := getTemplate(filepath.Join(util.GetBlogTemplatePath(), "blog_pages", "post.html"))
	if err != nil {
		logger.Error("获取模板失败: ", err.Error())
		return errors.New("获取模板失败: " + err.Error())
	}
	str := ReplaceBasePath(template, "")
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("<article><h1>%s</h1><div class=\"post-desc\"><time>%s</time>", article.Title, article.CreateAt))
	for i := range article.Tags {
		sb.WriteString(fmt.Sprintf("<span>%s</span>", article.Tags[i]))
	}
	sb.WriteString(fmt.Sprintf("</div><div class=\"post-content pswp-gallery\">%s</div></article>", article.HtmlText))
	str = strings.Replace(str, "#{{post_content}}", sb.String(), 1)

	return util.WriteFile([]byte(str), filepath.Join(service.GetSysConfig().SavePath, "hyakkei", article.Slug+".html"))
}

// 为自定义页生成对应 html 文件
func GenerateCustomPage(post *model.Post) error {
	pageGenerateLock.Lock()
	defer pageGenerateLock.Unlock()
	template, err := getTemplate(filepath.Join(util.GetBlogTemplatePath(), "blog_pages", "page.html"))
	if err != nil {
		logger.Error("获取模板失败: ", err.Error())
		return errors.New("获取模板失败: " + err.Error())
	}
	str := ReplaceBasePath(template, "..")
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("<article><h1>%s</h1><div class=\"post-desc\">", post.Title))
	sb.WriteString(fmt.Sprintf("</div><div class=\"post-content pswp-gallery\">%s</div></article>", post.HtmlText))
	str = strings.Replace(str, "#{{page_content}}", sb.String(), 1)

	return util.WriteFile([]byte(str), filepath.Join(service.GetSysConfig().SavePath, "hyakkei", post.Slug+".html"))
}

func getTemplate(path string) (result []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// 生成 Header 文件
func GenerateHeaderFile() error {
	template, err := getTemplate(filepath.Join(util.GetBlogTemplatePath(), "blog_pages", "header.html"))
	if err != nil {
		return errors.New("读取 header 模版文件数据失败: " + err.Error())
	}
	str := ReplaceBasePath(template, "")
	cfg := service.GetSysConfig()
	// 是否展示 github 链接
	githubStr := ""
	if cfg.IsShowGithub {
		githubStr = fmt.Sprintf(`<li><a href="https://github.com/%s" target="_blank">GitHub</a></li>`, cfg.GithubName)
	}
	str = strings.Replace(str, "#{{github}}", githubStr, 1)
	// 是否展示图书页面链接
	bookStr := ""
	if cfg.IsShowBook {
		bookStr = `<li><a href="books.html" target="_parent">Read</a></li>`
	}
	str = strings.Replace(str, "#{{book_page}}", bookStr, 1)
	str = strings.Replace(str, "#{{logo_filename}}", cfg.LogoName, 1)
	// 自定义页面菜单构建
	d, err := service.GetPagesData()
	if err != nil {
		return err
	}
	list := make([]model.Post, 0)
	_ = json.Unmarshal(d, &list)
	sb := strings.Builder{}
	for i := range list {
		if list[i].Status == model.Publish && list[i].Slug != "" {
			sb.WriteString(fmt.Sprintf(`<li><a href="%s.html" target="_parent">%s</a></li>`, list[i].Slug, list[i].Slug))
		}
	}
	str = strings.Replace(str, "#{{custom_page}}", sb.String(), 1)
	return util.WriteFile([]byte(str), filepath.Join(cfg.SavePath, "hyakkei", "header.html"))
}

// 生成 Footer 文件
func GenerateFooterFile() error {
	template, err := getTemplate(filepath.Join(util.GetBlogTemplatePath(), "blog_pages", "footer.html"))
	if err != nil {
		return errors.New("读取 footer 模版文件数据失败: " + err.Error())
	}
	cfg := service.GetSysConfig()
	return util.WriteFile([]byte(template), filepath.Join(cfg.SavePath, "hyakkei", "footer.html"))
}

const pswpStr = `<div class="pswp" tabindex="-1" role="dialog" aria-hidden="true">
        <!-- Background of PhotoSwipe.
         It's a separate element as animating opacity is faster than rgba(). -->
        <div class="pswp__bg"></div>
        <div class="pswp__scroll-wrap">
            <div class="pswp__container">
                <div class="pswp__item"></div>
                <div class="pswp__item"></div>
                <div class="pswp__item"></div>
            </div>
            <!-- Default (PhotoSwipeUI_Default) interface on top of sliding area. Can be changed. -->
            <div class="pswp__ui pswp__ui--hidden">
                <div class="pswp__top-bar">
                    <div class="pswp__counter"></div>
                    <button class="pswp__button pswp__button--close" title="Close (Esc)"></button>
                    <!-- <button class="pswp__button pswp__button--share" title="Share"></button> -->
                    <button class="pswp__button pswp__button--fs" title="Toggle fullscreen"></button>
                    <button class="pswp__button pswp__button--zoom" title="Zoom in/out"></button>
                    <!-- Preloader demo https://codepen.io/dimsemenov/pen/yyBWoR -->
                    <!-- element will get class pswp__preloader--active when preloader is running -->
                    <div class="pswp__preloader">
                        <div class="pswp__preloader__icn">
                            <div class="pswp__preloader__cut">
                                <div class="pswp__preloader__donut"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
                    <div class="pswp__share-tooltip"></div>
                </div>

                <button class="pswp__button pswp__button--arrow--left" title="Previous (arrow left)">
                </button>

                <button class="pswp__button pswp__button--arrow--right" title="Next (arrow right)">
                </button>

                <div class="pswp__caption">
                    <div class="pswp__caption__center"></div>
                </div>
            </div>
        </div>
    </div>`
