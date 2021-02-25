package ticker

import (
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/util"
	"path/filepath"
	"strconv"
	"time"
)

func RunTickerTask() {
	go autoUpdateVisitorCount()
}

func autoUpdateVisitorCount() {
	ticker := time.NewTicker(time.Minute * 4)
	for {
		select {
		case <-ticker.C:
			v := util.Cache.GetInt32(model.VisitorCount)
			path := filepath.Join(util.GetBlogDataPath(), "visitor_count.txt")
			if err := util.WriteFile([]byte(strconv.FormatInt(int64(v), 10)), path); err != nil {
				logger.Error("更新访客数量失败，", err.Error())
			}
		}
	}
}
