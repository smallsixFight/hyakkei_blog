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

var bookLock sync.Mutex

const BookFileName = "books.json"

func SaveBooks(bs []byte) error {
	bookLock.Lock()
	defer bookLock.Unlock()
	filename := filepath.Join(util.GetBlogDataPath(), BookFileName)
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

func FindBookAndIdx(id int64, list []model.Book) (idx int, book *model.Book) {
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

func BookIsExist(id int64, title string, list []model.Book) bool {
	for i := range list {
		if list[i].Title == title && list[i].Id != id {
			return true
		}
	}
	return false
}

func GetBooks() (d []byte, err error) {
	data := util.Cache.GetBytes("books")
	if data == nil {
		logger.Println("图书数据缓存不存在")
		d, err = getBooksData()
		if err != nil {
			return
		}
		data = d
	}
	return data, nil
}

func getBooksData() (data []byte, err error) {
	bookLock.Lock()
	defer bookLock.Unlock()
	filename := filepath.Join(util.GetBlogDataPath(), BookFileName)
	f, err := os.Open(filename)
	if err != nil {
		logger.Println(err.Error())
		err = fmt.Errorf("文件[%s]不存在或被意外删除，请检查。", filename)
		return
	}
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	if err == nil {
		util.Cache.Set("books", data, 0)
	}
	return
}
