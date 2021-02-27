package file_generator

import (
	"encoding/json"
	"fmt"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"os"
	"path/filepath"
)

var (
	BlogPagePath   = filepath.Join(util.GetBlogTemplatePath(), "blog_pages")
	BlogFileList   = []string{"index.html", "post.html", "page.html", "friends.html", "books.html", "header.html"}
	BlogBSPagePath = filepath.Join(util.GetBlogTemplatePath(), "backstage_pages")
	BlogBsFileList = []string{"articles.html", "articles_edit.html", "friends.html", "login.html",
		"pages.html", "pages_edit.html", "preview.html", "setting.html", "tags.html", "dashboard.html"}
)

func GenerateBlogFile() error {
	if err := util.CopyDir(filepath.Join(util.GetBlogTemplatePath(), "assets"),
		filepath.Join(service.GetSysConfig().SavePath, "hyakkei")); err != nil {
		return err
	}
	// 生成前端页面
	if err := GenerateIndexPage(); err != nil {
		return err
	}
	if err := GenerateFriendLinksPage(); err != nil {
		return err
	}
	if err := GenerateBookPage(); err != nil {
		return err
	}
	if err := GenerateHeaderFile(); err != nil {
		return err
	}
	// 创建文章页
	posts := make([]model.Post, 0)
	d, err := service.GetArticlesData()
	if err != nil {
		return err
	}
	_ = json.Unmarshal(d, &posts)
	for i := range posts {
		if err := GenerateArticlePage(&posts[i]); err != nil {
			return err
		}
	}
	// 创建自定义页
	d, err = service.GetPagesData()
	if err != nil {
		return err
	}
	posts = make([]model.Post, 0)
	_ = json.Unmarshal(d, &posts)
	savePath := filepath.Join(service.GetSysConfig().SavePath, "hyakkei", "custom_page")
	if !util.FileIsExist(savePath) {
		if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
			return err
		}
	}
	for i := range posts {
		if err := GenerateCustomPage(&posts[i]); err != nil {
			return err
		}
	}
	// 生成后台管理页面
	savePath = filepath.Join(service.GetSysConfig().SavePath, "hyakkei", "bs")
	if !util.FileIsExist(savePath) {
		if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
			return err
		}
	}
	for i := range BlogBsFileList {
		template, err := getTemplate(filepath.Join(BlogBSPagePath, BlogBsFileList[i]))
		if err != nil {
			return fmt.Errorf("读取 %s 模版文件数据失败: %s", BlogBsFileList[i], err.Error())
		}
		str := ReplaceBasePath(template)
		if err := util.WriteFile([]byte(str), filepath.Join(savePath, BlogBsFileList[i])); err != nil {
			return err
		}
	}
	return nil
}

func GenerateBlogConfig(sysCfg *model.SysSetting) error {
	if err := util.SaveInitConfig(&util.InitConfig{
		JsPath:   "assets/js",
		ImgPath:  "assets/img",
		CssPath:  "assets/css",
		FontPath: "assets/font",
	}); err != nil {
		return err
	}
	return service.GenerateSysSetting(sysCfg)
}

func GenerateDataFile(arr []string) error {
	for i := range arr {
		path := filepath.Join(util.GetBlogDataPath(), arr[i])
		if util.FileIsExist(path) {
			continue
		}
		fmt.Printf("正在创建[%s]\n" + path)
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("创建失败: %s", err.Error())
		}
		if _, err := f.WriteString("[]"); err != nil {
			f.Close()
			return fmt.Errorf("写入初始数据失败: %s", err.Error())
		}
		f.Close()
	}
	return nil
}
