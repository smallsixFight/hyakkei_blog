package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"net/http"
)

func GetDashboardInfo(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	info := &model.Dashboard{InitTime: util.GetSysConfig().InitTime}

	articles := make([]model.Post, 0)
	pages := make([]model.Post, 0)
	data, err := service.GetArticlesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	_ = json.Unmarshal(data, &articles)
	data, err = service.GetPagesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	_ = json.Unmarshal(data, &pages)
	for i := range articles {
		if articles[i].Status == model.Publish {
			info.ArticleInfo.PublishCount++
			info.ArticleInfo.LastPublish = articles[i].Title
		} else {
			info.ArticleInfo.DraftCount++
		}
	}
	info.ArticleInfo.LastAdd = articles[len(articles)-1].Title

	for i := range pages {
		if pages[i].Status == model.Publish {
			info.PageInfo.PublishCount++
			info.PageInfo.LastPublish = pages[i].Title
		} else {
			info.PageInfo.DraftCount++
		}
	}
	info.PageInfo.LastAdd = pages[len(pages)-1].Title

	info.FriendLinkCount = util.Cache.GetInt32(model.FriendLinkCount)
	info.BookCount = util.Cache.GetInt32(model.BookCount)
	info.VisitorCount = util.Cache.GetInt32(model.VisitorCount)
	reply.SetData(info).SetSuccess(true).SetMessage("ok")
}
