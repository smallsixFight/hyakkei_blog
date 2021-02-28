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

func GetTagList(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	page := util.Str2Int(ctx.Query("page"))
	typ := ctx.Query("typ")
	if page < 1 {
		page = 1
	}
	list, err := service.GetTags()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.Tag, 0)
	reply.SetData(&data)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析标签数据失败")
		return
	}
	if typ != "all" && (page-1)*10 < len(data) {
		last := page * 10
		if last > len(data) {
			last = len(data)
		}
		reply.SetData(data[(page-1)*10 : last])
	}
	reply.SetSuccess(true).SetTotal(int64(len(data))).SetMessage("ok")
}

func DeleteTag(ctx *gin.Context) {
	reply := model.Reply{Message: "标签不存在"}
	defer ctx.JSON(http.StatusOK, &reply)
	id := util.Str2Int64(ctx.Query("id"))
	if id == 0 {
		reply.SetMessage("id 不能为空")
		return
	}
	list, err := service.GetTags()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.Tag, 0)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析标签数据失败")
		return
	}
	idx, _ := service.FindTagAndIdx(id, data)
	if idx == -1 {
		return
	}
	data = append(data[:idx], data[idx+1:]...)
	bs, _ := json.Marshal(data)
	if err := service.SaveTags(bs); err != nil {
		reply.SetMessage("标签列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set("tags", bs, 0)
	reply.SetMessage("删除成功").SetSuccess(true)
}

func SaveTag(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	tag := model.Tag{}
	if err := ctx.BindJSON(&tag); err != nil {
		logger.Println("提交参数错误: ", err.Error())
		reply.Message = "提交参数错误"
		return
	}
	if err := handleTagParam(&tag); err != nil {
		reply.SetMessage(err.Error())
		return
	}
	tag.ModifyAt = time.Now().Format("2006-01-02 15:04:05")
	list, err := service.GetTags()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	tags := make([]model.Tag, 0)
	if err := json.Unmarshal(list, &tags); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析标签列表数据失败")
		return
	}
	// 判断标签是否存在
	if service.TagIsExist(tag.Id, tag.Name, tags) {
		reply.SetMessage("该标签已存在，请修改")
		return
	}
	// 保存
	// id 不为 0 并且标签存在则只需要更新，不用新增
	var isUpdate bool
	if tag.Id != 0 {
		if idx, v := service.FindTagAndIdx(tag.Id, tags); idx != -1 {
			tag.CreateAt = v.CreateAt
			tags[idx] = tag
			isUpdate = true
		} else {
			reply.SetMessage("更新失败: 标签不存在")
			return
		}
	}
	if !isUpdate {
		tag.Id = util.GetId()
		tag.CreateAt = tag.ModifyAt
		tags = append(tags, tag)
	}
	bs, _ := json.Marshal(&tags)
	if err := service.SaveTags(bs); err != nil {
		reply.SetMessage("标签列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set("tags", bs, 0)

	reply.SetSuccess(true).SetMessage("保存成功")
}

func handleTagParam(tag *model.Tag) error {
	tag.Name = strings.TrimSpace(tag.Name)
	tag.Description = strings.TrimSpace(tag.Description)
	if tag.Name == "" {
		return errors.New("名称不能为空")
	}
	return nil
}
