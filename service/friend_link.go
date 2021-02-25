package service

import (
	"fmt"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

func FindFriendLinkAndIdx(id int64, list []model.FriendLink) (idx int, link *model.FriendLink) {
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

func SaveFriendLinks(bs []byte) error {
	filename := filepath.Join(util.GetBlogDataPath(), "friend_links.json")
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logger.Println(err.Error())
		return fmt.Errorf("友链文件[%s]不存在或被意外删除，请检查。", filename)
	}
	defer f.Close()
	if _, err := f.Write(bs); err != nil {
		return err
	}
	return nil
}

func GetFriendLinks() (d []byte, err error) {
	data := util.Cache.GetBytes("friend_links")
	if data == nil {
		logger.Println("友链数据缓存不存在")
		d, err = getFriendLinksData()
		if err != nil {
			return
		}
		data = d
	}
	return data, nil
}

func FriendLinkIsExist(id int64, url string, list []model.FriendLink) bool {
	for i := range list {
		if list[i].Url == url && list[i].Id != id {
			return true
		}
	}
	return false
}

func getFriendLinksData() (data []byte, err error) {
	filename := filepath.Join(util.GetBlogDataPath(), "friend_links.json")
	f, err := os.Open(filename)
	if err != nil {
		logger.Println(err.Error())
		err = fmt.Errorf("友链文件[%s]不存在或被意外删除，请检查。", filename)
		return
	}
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	if err == nil {
		util.Cache.Set("friend_links", data, 0)
	}
	return
}
