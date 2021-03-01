package blog_install

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"github.com/smallsixFight/hyakkei_blog/util/file_generator"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Init() error {
	if err := file_generator.GenerateBlogConfig(setSysConfig()); err != nil {
		return err
	}
	// 创建数据存储文件
	dataNameList := []string{"articles.json", "books.json", "friend_links.json", "pages.json", "tags.json"}
	if err := file_generator.GenerateDataFile(dataNameList); err != nil {
		return err
	}
	if err := file_generator.GenerateBlogFile(); err != nil {
		return err
	}
	// 创建访客记录文件，并初始化缓存
	path := filepath.Join(util.GetBlogDataPath(), "visitor_count.txt")
	if !util.FileIsExist(path) {
		if err := util.WriteFile([]byte("0"), path); err != nil {
			return errors.New("生成访客记录文件失败: " + err.Error())
		}
	}
	return nil
}

func IsNeedInstall() bool {
	// 检查模板文件是否完整
	log.Println("正在校验模板文件是否完整...")
	cssPath := filepath.Join(util.GetBlogTemplatePath(), "assets", "css")
	cssFileList := []string{"default-skin.css", "default-skin.svg", "highlight.css", "main.css", "photoswipe.css"}
	for i := 0; i < len(cssFileList); i++ {
		if !util.FileIsExist(filepath.Join(cssPath, cssFileList[i])) {
			log.Fatalf("文件 %s 缺失，请检查！\n", cssFileList[i])
		}
	}

	editIconPath := filepath.Join(util.GetBlogTemplatePath(), "assets", "img", "icons", "editor")
	editIconFileList := []string{"bold.svg", "bulletedlist.svg", "code_block.svg", "hr.svg",
		"image.svg", "italic.svg", "link.svg", "numberedlist.svg", "quote.svg",
		"strikethrough.svg", "tag_code.svg", "underline.svg"}
	for i := 0; i < len(editIconFileList); i++ {
		if !util.FileIsExist(filepath.Join(editIconPath, editIconFileList[i])) {
			log.Fatalf("文件 %s 缺失，请检查！\n", editIconFileList[i])
		}
	}

	jsPath := filepath.Join(util.GetBlogTemplatePath(), "assets", "js")
	jsFileList := []string{"gallery-init.js", "highlight.js", "jquery.min.js", "main.js",
		"photoswipe.min.js", "photoswipe-ui-default.min.js"}
	for i := 0; i < len(jsFileList); i++ {
		if !util.FileIsExist(filepath.Join(jsPath, jsFileList[i])) {
			log.Fatalf("文件 %s 缺失，请检查！\n", jsFileList[i])
		}
	}
	// 前台页面
	for i := 0; i < len(file_generator.BlogFileList); i++ {
		if !util.FileIsExist(filepath.Join(file_generator.BlogPagePath, file_generator.BlogFileList[i])) {
			log.Fatalf("文件 %s 缺失，请检查！\n", file_generator.BlogFileList[i])
		}
	}
	// 后台页面
	for i := 0; i < len(file_generator.BlogBsFileList); i++ {
		if !util.FileIsExist(filepath.Join(file_generator.BlogBSPagePath, file_generator.BlogBsFileList[i])) {
			log.Fatalf("文件 %s 缺失，请检查！\n", file_generator.BlogBsFileList[i])
		}
	}
	log.Println("模板文件检测完成，无缺失文件。")
	log.Println("正在检查配置文件是否缺失...")
	// 检查配置文件是否创建
	if !util.FileIsExist(filepath.Join(util.GetBlogConfigPath(), util.InitConfigName)) ||
		!util.FileIsExist(filepath.Join(util.GetBlogConfigPath(), service.SysConfigName)) {
		log.Println("配置文件缺失，将重新生成。")
		return true
	}
	if !util.FileIsExist(filepath.Join(util.GetBlogDataPath(), "visitor_count.txt")) {
		log.Println("访客记录文件缺失，将重新生成。")
		return true
	}
	log.Println("检查完成，配置文件无缺失。")
	return false
}

func SetRunCache() error {
	// 初始化系统设置
	if err := service.InitSysConfig(); err != nil {
		return err
	}
	// 读取访客记录文件，并缓存数据
	visitorRecordPath := filepath.Join(util.GetBlogDataPath(), "visitor_count.txt")
	data, err := ioutil.ReadFile(visitorRecordPath)
	if err != nil {
		return errors.New("读取访客记录文件失败， " + err.Error())
	}
	util.Cache.Set(model.VisitorCount, util.Str2Int32(string(data)), 0)
	// 初始化友链数、书籍数缓存
	friends := make([]model.FriendLink, 0)
	books := make([]model.Book, 0)
	data, err = service.GetFriendLinks()
	if err != nil {
		return err
	}
	_ = json.Unmarshal(data, &friends)
	data, err = service.GetBooks()
	if err != nil {
		return err
	}
	_ = json.Unmarshal(data, &books)

	util.Cache.Set(model.FriendLinkCount, int32(len(friends)), 0)
	util.Cache.Set(model.BookCount, int32(len(books)), 0)

	return nil
}

func setSysConfig() *model.SysSetting {
	cfg := &model.SysSetting{}
	for {
		if cfg.BlogName == "" {
			v := strings.TrimSpace(getInput("博客名称"))
			if v == "" {
				fmt.Println("博客名称不能为空")
				continue
			}
			cfg.BlogName = v
		}
		if cfg.Username == "" {
			v := strings.TrimSpace(getInput("用户名"))
			if v == "" {
				fmt.Println("用户名不能为空")
				continue
			}
			cfg.Username = v
		}
		if cfg.Password == "" {
			v := strings.TrimSpace(getInput("登录密码"))
			if v == "" {
				fmt.Println("登录密码不能为空")
				continue
			}
			cfg.Password = v
		}
		if cfg.Salt == "" {
			v := strings.TrimSpace(getInput("加密 salt(default: hyakkei)"))
			if v == "" {
				v = "hyakkei"
			}
			cfg.Salt = v
		}
		if cfg.TokenSecret == "" {
			v := strings.TrimSpace(getInput("token secret"))
			if v == "" {
				fmt.Println("token secret 必须设置")
				continue
			}
			cfg.TokenSecret = v
		}
		showGithub := strings.ToLower(strings.TrimSpace(getInput("是否显示 Github 链接(y/n, default: y)")))
		if showGithub != "n" {
			cfg.IsShowGithub = true
			if cfg.GithubName == "" {
				v := strings.TrimSpace(getInput("github 名称(e.g. smallsixFight)"))
				cfg.GithubName = v
			}
		}
		showBook := strings.ToLower(strings.TrimSpace(getInput("是否显示图书页(y/n, default: n)")))
		if showBook == "y" {
			cfg.IsShowBook = true
		}
		if cfg.SavePath == "" {
			v := strings.TrimSpace(getInput("文件存储路径(default: 当前路径)"))
			if v != "" {
				if !util.IsDir(v) && !filepath.IsAbs(v) {
					fmt.Println("请输入正确的路径（需绝对路径）")
					continue
				}
				if !util.FileIsExist(v) {
					if err := os.MkdirAll(v, os.ModePerm); err != nil {
						logger.Fatal("创建保存路径失败，" + err.Error())
						return nil
					}
				}
				cfg.SavePath = v
			} else {
				cfg.SavePath = util.GetAbsPath()
			}
		}
		break
	}
	cfg.LogoName = "logo.png"
	cfg.InitTime = time.Now().Unix()
	return cfg
}

func getInput(title string) string {
	fmt.Printf("%s: ", title)
	in := bufio.NewReader(os.Stdin)
	bs, _, _ := in.ReadLine()
	return string(bs)
}
