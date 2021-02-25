package service

import (
	"fmt"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	ArticleFileName = "articles.json"
	PageFileName    = "pages.json"
)

var postLock sync.Mutex

func GetPagesData() (d []byte, err error) {
	data := util.Cache.GetBytes("pages")
	if data == nil {
		logger.Println("页面缓存不存在")
		d, err = getPost(PageFileName, "pages")
		if err != nil {
			return
		}
		data = d
	}
	return data, nil
}

func GetArticlesData() (d []byte, err error) {
	data := util.Cache.GetBytes("articles")
	if data == nil {
		logger.Println("文章缓存不存在")
		d, err = getPost(ArticleFileName, "articles")
		if err != nil {
			return
		}
		data = d
	}
	return data, nil
}

func getPost(postName, cacheKey string) (data []byte, err error) {
	postLock.Lock()
	defer postLock.Unlock()
	filename := filepath.Join(util.GetBlogDataPath(), postName)
	f, err := os.Open(filename)
	if err != nil {
		logger.Println(err.Error())
		err = fmt.Errorf("文件[%s]不存在或被意外删除，请检查。", filename)
		return
	}
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	if err == nil {
		util.Cache.Set(cacheKey, data, 0)
	}
	return
}

func savePosts(bs []byte, postName string) error {
	postLock.Lock()
	defer postLock.Unlock()
	filename := filepath.Join(util.GetBlogDataPath(), postName)
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logger.Println(err.Error())
		return fmt.Errorf("文件[%s]不存在或被意外删除，请检查。", filename)
	}
	defer f.Close()
	if _, err := f.Write(bs); err != nil {
		return err
	}
	return nil
}

func SaveArticles(bs []byte) error {
	return savePosts(bs, ArticleFileName)
}

func SavePages(bs []byte) error {
	return savePosts(bs, PageFileName)
}

func PostIsExist(id int64, filename string, list []model.Post) bool {
	for i := range list {
		if list[i].Slug == filename && list[i].Id != id {
			return true
		}
	}
	return false
}

func FindPost(id int64, list []model.Post) *model.Post {
	_, d := FindPostAndIdx(id, list)
	return d
}

func FindPostAndIdx(id int64, list []model.Post) (idx int, page *model.Post) {
	low, high := 0, len(list)-1
	for low <= high {
		mid := (low + high) / 2
		if list[mid].Id == id {
			return mid, &list[mid]
		} else if list[mid].Id > id {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	idx = -1
	return
}
