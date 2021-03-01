package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"net/http"
)

func FetchArticleList(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	page := util.Str2Int(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	bs, err := service.GetArticlesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	list := make([]model.BasePostInfo, 0)
	_ = json.Unmarshal(bs, &list)
	low, high := 0, len(list)-1
	for low < high {
		list[low], list[high] = list[high], list[low]
		low++
		high--
	}
	total := 0
	addCount := 0
	skip := (page - 1) * 10
	data := make([]model.BasePostInfo, 0)
	for i := range list {
		if list[i].Status == model.Publish {
			total++
			if skip > 0 {
				skip--
				continue
			}
			if addCount < 10 {
				data = append(data, list[i])
				addCount++
			}
		}
	}
	reply.SetSuccess(true).SetTotal(int64(total)).SetMessage("ok").SetData(data)
}

func GetBooks(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	bs, err := service.GetBooks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	list := make([]model.Book, 0)
	_ = json.Unmarshal(bs, &list)
	data := struct {
		WishList []model.Book `json:"wish_list"`
		ReadList []model.Book `json:"read_list"`
	}{make([]model.Book, 0), make([]model.Book, 0)}
	for i := range list {
		if list[i].Status == model.Wish {
			data.WishList = append(data.WishList, list[i])
		} else {
			data.ReadList = append(data.ReadList, list[i])
		}
	}
	reply.SetData(data).SetSuccess(true).SetMessage("ok")
}

func GetFriends(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	bs, err := service.GetFriendLinks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	list := make([]model.FriendLink, 0)
	_ = json.Unmarshal(bs, &list)
	reply.SetData(list).SetSuccess(true).SetMessage("ok")
}
