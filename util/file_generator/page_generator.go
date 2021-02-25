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
	str := util.ReplaceBasePath(template)
	return util.WriteFile([]byte(str), filepath.Join(dest, pageName))
}

// 生成初始页
func GenerateIndexPage() error {
	return GenerateBasePage(filepath.Join(util.GetBlogTemplatePath(), "blog_pages"), filepath.Join(util.GetSysConfig().SavePath, "hyakkei"), "index.html")
}

// 生成书籍页
func GenerateBookPage() error {
	return GenerateBasePage(filepath.Join(util.GetBlogTemplatePath(), "blog_pages"), filepath.Join(util.GetSysConfig().SavePath, "hyakkei"), "books.html")
}

// 生成友链页
func GenerateFriendLinksPage() error {
	return GenerateBasePage(filepath.Join(util.GetBlogTemplatePath(), "blog_pages"), filepath.Join(util.GetSysConfig().SavePath, "hyakkei"), "friends.html")
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
	str := util.ReplaceBasePath(template)
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("<article><h1>%s</h1><div class=\"post-desc\"><time>%s</time>", article.Title, article.CreateAt))
	for i := range article.Tags {
		sb.WriteString(fmt.Sprintf("<span>%s</span>", article.Tags[i]))
	}
	sb.WriteString(fmt.Sprintf("</div><div class=\"post-content pswp-gallery\">%s</div></article>", article.HtmlText))
	str = strings.Replace(str, "#{{post_content}}", sb.String(), 1)

	return util.WriteFile([]byte(str), filepath.Join(util.GetSysConfig().SavePath, "hyakkei", article.Slug+".html"))
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
	str := util.ReplaceBasePath(template)
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("<article><h1>%s</h1><div class=\"post-desc\">", post.Title))
	sb.WriteString(fmt.Sprintf("</div><div class=\"pswp-gallery\">%s</div></article>", post.HtmlText))
	str = strings.Replace(str, "#{{page_content}}", sb.String(), 1)

	return util.WriteFile([]byte(str), filepath.Join(util.GetSysConfig().SavePath, "hyakkei", "custom_page", post.Slug+".html"))
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
	str := util.ReplaceBasePath(template)
	cfg := util.GetSysConfig()
	// 是否展示 github 链接
	githubStr := ""
	if cfg.IsShowGithub {
		githubStr = fmt.Sprintf(`<li><a href="https://github.com/%s" target="_blank">GitHub</a></li>`, cfg.GithubName)
	}
	str = strings.Replace(str, "#{{github}}", githubStr, 1)
	// 是否展示图书页面链接
	bookStr := ""
	if cfg.IsShowBook {
		bookStr = `<li><a href="books.html" target="_blank">Read</a></li>`
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
			sb.WriteString(fmt.Sprintf(`<li><a href="%s.html" target="_parent">%s</a></li>`, filepath.Join("custom_page", list[i].Slug), list[i].Slug))
		}
	}
	str = strings.Replace(str, "#{{custom_page}}", sb.String(), 1)
	return util.WriteFile([]byte(str), filepath.Join(cfg.SavePath, "hyakkei", "header.html"))
}
