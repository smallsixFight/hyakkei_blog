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

func GetArticleList(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	page := util.Str2Int(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	list, err := service.GetArticlesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.BasePostInfo, 0)
	reply.SetData(data)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析文章数据失败")
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

func GetArticleDetail(ctx *gin.Context) {
	reply := model.Reply{Message: "文章不存在"}
	defer ctx.JSON(http.StatusOK, &reply)
	id := util.Str2Int64(ctx.Param("id"))
	if id == 0 {
		reply.SetMessage("id 不能为空")
		return
	}
	list, err := service.GetArticlesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.Post, 0)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析文章数据失败")
		return
	}
	article := service.FindPost(id, data)
	if article != nil {
		reply.SetData(article).SetSuccess(true).SetMessage("ok")
	}
}

func SaveArticle(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	article := model.Post{}
	if err := ctx.BindJSON(&article); err != nil {
		logger.Warn("保存文章提交参数错误: ", err.Error())
		reply.Message = "提交参数错误"
		return
	}
	if err := handlePostParam(&article, "article"); err != nil {
		reply.SetMessage(err.Error())
		return
	}
	article.ModifyAt = time.Now().Format("2006-01-02 15:04:05")
	filename := article.Slug
	if filename == "" {
		article.Slug = article.Title
		filename = article.Title
	}
	list, err := service.GetArticlesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	articles := make([]model.Post, 0)
	if err := json.Unmarshal(list, &articles); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析文章列表数据失败")
		return
	}
	// 判断文章标题是否存在
	if service.PostIsExist(article.Id, filename, articles) {
		reply.SetMessage("文章标题/slug 已存在，请修改")
		return
	}
	var b bytes.Buffer
	if err := util.GetMDHandle().Convert([]byte(article.MarkdownText), &b); err != nil {
		reply.SetMessage("文本转换失败: " + err.Error())
		return
	}
	article.HtmlText = b.String()
	// 保存文章
	// id 不为 0 并且文章存在则只需要更新，不用新增
	var isUpdate bool
	var preName string
	if article.Id != 0 {
		if idx, v := service.FindPostAndIdx(article.Id, articles); idx != -1 {
			article.CreateAt = v.CreateAt
			articles[idx] = article
			isUpdate = true
			// 文章之前发布过，则删除之前生成的文件
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
			reply.SetMessage("更新失败: 文章不存在")
			return
		}
	}
	if !isUpdate {
		article.Id = util.GetId()
		article.CreateAt = article.ModifyAt
		articles = append(articles, article)
	}
	bs, _ := json.Marshal(&articles)
	if err := service.SaveArticles(bs); err != nil {
		reply.SetMessage("文章列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set("articles", bs, 0)

	// 发布的文章创建对应的文件
	if article.Status == model.Publish {
		if err := file_generator.GenerateArticlePage(&article); err != nil {
			reply.SetMessage(err.Error())
			return
		}
	}
	reply.SetSuccess(true).SetMessage("保存成功")
}

func DeleteArticle(ctx *gin.Context) {
	reply := model.Reply{Message: "文章不存在"}
	defer ctx.JSON(http.StatusOK, &reply)
	id := util.Str2Int64(ctx.Query("id"))
	if id == 0 {
		reply.SetMessage("id 不能为空")
		return
	}
	list, err := service.GetArticlesData()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.Post, 0)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析文章数据失败")
		return
	}
	idx, val := service.FindPostAndIdx(id, data)
	if idx == -1 {
		return
	}
	if val.Status == model.Publish {
		filename := filepath.Join(util.GetAbsPath(), "hyakkei", getPostName(val)+".html")
		if err := os.Remove(filename); err != nil {
			logger.Warnf("[%s]删除失败: %s", filename, err.Error())
			return
		}
	}
	data = append(data[:idx], data[idx+1:]...)
	bs, _ := json.Marshal(data)
	if err := service.SaveArticles(bs); err != nil {
		reply.SetMessage("文章列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set("articles", bs, 0)
	reply.SetMessage("删除成功").SetSuccess(true)
}
