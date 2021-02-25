package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"net/http"
	"strings"
	"time"
)

func GetFriendLinkList(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	page := util.Str2Int(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	list, err := service.GetFriendLinks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.FriendLink, 0)
	reply.SetData(data)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析友链数据失败")
		return
	}
	if (page-1)*10 < len(data) {
		last := page * 10
		if last > len(data) {
			last = len(data)
		}
		reply.SetData(data[(page-1)*10 : last])
	}
	reply.SetSuccess(true).SetTotal(int64(len(data))).SetMessage("ok")
}

func DeleteFriendLink(ctx *gin.Context) {
	reply := model.Reply{Message: "友链不存在"}
	defer ctx.JSON(http.StatusOK, &reply)
	id := util.Str2Int64(ctx.Query("id"))
	if id == 0 {
		reply.SetMessage("id 不能为空")
		return
	}
	list, err := service.GetFriendLinks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.FriendLink, 0)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析友链数据失败")
		return
	}
	idx, _ := service.FindFriendLinkAndIdx(id, data)
	if idx == -1 {
		return
	}
	data = append(data[:idx], data[idx+1:]...)
	bs, _ := json.Marshal(data)
	if err := service.SaveFriendLinks(bs); err != nil {
		reply.SetMessage("友链列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set(model.FriendLinkCount, len(data), 0)
	util.Cache.Set("friend_links", bs, 0)
	reply.SetMessage("删除成功").SetSuccess(true)
}

func SaveFriendLink(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	link := model.FriendLink{}
	if err := ctx.BindJSON(&link); err != nil {
		logger.Println("提交参数错误: ", err.Error())
		reply.Message = "提交参数错误"
		return
	}
	if err := handleLinkParam(&link); err != nil {
		reply.SetMessage(err.Error())
		return
	}
	link.ModifyAt = time.Now().Format("2006-01-02 15:04:05")
	list, err := service.GetFriendLinks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	links := make([]model.FriendLink, 0)
	if err := json.Unmarshal(list, &links); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析友链列表数据失败")
		return
	}
	// 判断是否存在
	if service.FriendLinkIsExist(link.Id, link.Url, links) {
		reply.SetMessage("该友链已存在，请修改")
		return
	}
	// 保存
	// id 不为 0 并且友链存在则只需要更新，不用新增
	var isUpdate bool
	if link.Id != 0 {
		if idx, v := service.FindFriendLinkAndIdx(link.Id, links); idx != -1 {
			link.CreateAt = v.CreateAt
			links[idx] = link
			isUpdate = true
		} else {
			reply.SetMessage("更新失败: 友链不存在")
			return
		}
	}
	if !isUpdate {
		link.Id = util.GetId()
		link.CreateAt = link.ModifyAt
		links = append(links, link)
	}
	bs, _ := json.Marshal(&links)
	if err := service.SaveFriendLinks(bs); err != nil {
		reply.SetMessage("友链列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set(model.FriendLinkCount, len(links), 0)
	util.Cache.Set("friend_links", bs, 0)

	reply.SetSuccess(true).SetMessage("保存成功")
}

func handleLinkParam(link *model.FriendLink) error {
	link.Name = strings.TrimSpace(link.Name)
	link.Url = strings.TrimSpace(link.Url)
	link.AvatarUrl = strings.TrimSpace(link.AvatarUrl)
	link.Description = strings.TrimSpace(link.Description)
	if link.Name == "" {
		return errors.New("名称不能为空")
	} else if link.Url == "" {
		return errors.New("链接地址不能为空")
	} else if !util.IsUrl(link.Url) {
		return errors.New("请输入有效的链接地址")
	} else if link.AvatarUrl != "" && !util.IsUrl(link.AvatarUrl) {
		return errors.New("请输入有效的头像地址")
	}
	return nil
}
