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

var tagLock sync.Mutex

func TagIsExist(id int64, name string, list []model.Tag) bool {
	for i := range list {
		if list[i].Name == name && list[i].Id != id {
			return true
		}
	}
	return false
}

func SaveTags(bs []byte) error {
	tagLock.Lock()
	defer tagLock.Unlock()
	filename := filepath.Join(util.GetBlogDataPath(), "tags.json")
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logger.Println(err.Error())
		return fmt.Errorf("标签文件[%s]不存在或被意外删除，请检查。", filename)
	}
	defer f.Close()
	if _, err := f.Write(bs); err != nil {
		return err
	}
	return nil
}

func FindTagAndIdx(id int64, list []model.Tag) (idx int, tag *model.Tag) {
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

func GetTags() (d []byte, err error) {
	data := util.Cache.GetBytes("tags")
	if data == nil {
		logger.Println("标签数据缓存不存在")
		d, err = getTagsData()
		if err != nil {
			return
		}
		data = d
	}
	return data, nil
}

func getTagsData() (data []byte, err error) {
	tagLock.Lock()
	defer tagLock.Unlock()
	filename := filepath.Join(util.GetBlogDataPath(), "tags.json")
	f, err := os.Open(filename)
	if err != nil {
		logger.Println(err.Error())
		err = fmt.Errorf("标签文件[%s]不存在或被意外删除，请检查。", filename)
		return
	}
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	if err == nil {
		util.Cache.Set("tags", data, 0)
	}
	return
}
