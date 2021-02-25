package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"github.com/smallsixFight/hyakkei_blog/util/file_generator"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func GetPageList(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	page := util.Str2Int(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	list, err := service.GetPagesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.BasePostInfo, 0)
	reply.SetData(data)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析自定义页面数据失败")
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

func GetPageDetail(ctx *gin.Context) {
	reply := model.Reply{Message: "该页面不存在"}
	defer ctx.JSON(http.StatusOK, &reply)
	id := util.Str2Int64(ctx.Param("id"))
	if id == 0 {
		reply.SetMessage("id 不能为空")
		return
	}
	list, err := service.GetPagesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.Post, 0)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析自定义页数据失败")
		return
	}
	pageInfo := service.FindPost(id, data)
	if pageInfo != nil {
		reply.SetData(pageInfo).SetSuccess(true).SetMessage("ok")
	}
}

func SavePage(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	pageInfo := model.Post{}
	if err := ctx.BindJSON(&pageInfo); err != nil {
		logger.Println("提交参数错误: ", err.Error())
		reply.Message = "提交参数错误"
		return
	}
	if err := handlePostParam(&pageInfo, "page"); err != nil {
		reply.SetMessage(err.Error())
		return
	}
	pageInfo.ModifyAt = time.Now().Format("2006-01-02 15:04:05")
	filename := pageInfo.Slug
	if filename == "" {
		pageInfo.Slug = pageInfo.Title
		filename = pageInfo.Title
	}
	list, err := service.GetPagesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	pages := make([]model.Post, 0)
	if err := json.Unmarshal(list, &pages); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析自定义页列表数据失败")
		return
	}
	// 判断自定义页标题是否存在
	if service.PostIsExist(pageInfo.Id, filename, pages) {
		reply.SetMessage("自定义页标题/slug 已存在，请修改")
		return
	}
	var b bytes.Buffer
	if err := util.GetMDHandle().Convert([]byte(pageInfo.MarkdownText), &b); err != nil {
		reply.SetMessage("文本转换失败: " + err.Error())
		return
	}
	pageInfo.HtmlText = b.String()
	// 保存
	// id 不为 0 并且自定义页存在则只需要更新，不用新增
	var isUpdate bool
	var preName string
	if pageInfo.Id != 0 {
		if idx, v := service.FindPostAndIdx(pageInfo.Id, pages); idx != -1 {
			pageInfo.CreateAt = v.CreateAt
			pages[idx] = pageInfo
			isUpdate = true
			// 之前发布过，则删除之前生成的文件
			if v.Status == model.Publish {
				preName = v.Slug
				if preName == "" {
					preName = v.Title
				}
				filename := filepath.Join(util.GetAbsPath(), "hyakkei", preName+".html")
				if err := os.Remove(filename); err != nil {
					logger.Warnf("[%s]删除失败: %s", filename, err.Error())
				}
			}
		} else {
			reply.SetMessage("更新失败: 该页面不存在")
			return
		}
	}
	if !isUpdate {
		pageInfo.Id = util.GetId()
		pageInfo.CreateAt = pageInfo.ModifyAt
		pages = append(pages, pageInfo)
	}
	bs, _ := json.Marshal(&pages)
	if err := service.SavePages(bs); err != nil {
		reply.SetMessage("自定义页列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set("pages", bs, 0)

	// 发布的创建对应的文件
	if pageInfo.Status == model.Publish {
		if err := file_generator.GenerateCustomPage(&pageInfo); err != nil {
			reply.SetMessage(err.Error())
			return
		}
	}
	if err = file_generator.GenerateHeaderFile(); err != nil {
		logger.Warn("更新自定义页重新生成 header 文件失败: " + err.Error())
	}
	reply.SetSuccess(true).SetMessage("保存成功")
}

func DeletePage(ctx *gin.Context) {
	reply := model.Reply{Message: "删除失败，页面不存在"}
	defer ctx.JSON(http.StatusOK, &reply)
	id := util.Str2Int64(ctx.Query("id"))
	if id == 0 {
		reply.SetMessage("id 不能为空")
		return
	}
	list, err := service.GetPagesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.Post, 0)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析自定义页数据失败")
		return
	}
	idx, val := service.FindPostAndIdx(id, data)
	if idx == -1 {
		return
	}
	if val.Status == model.Publish {
		filename := filepath.Join(util.GetAbsPath(), "hyakkei", "custom_page", getPostName(val)+".html")
		if err := os.Remove(filename); err != nil {
			logger.Warnf("[%s]删除失败: %s", filename, err.Error())
			return
		}
	}
	data = append(data[:idx], data[idx+1:]...)
	bs, _ := json.Marshal(data)
	if err := service.SavePages(bs); err != nil {
		reply.SetMessage("自定义页列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set("pages", bs, 0)
	reply.SetMessage("删除成功").SetSuccess(true)
}
