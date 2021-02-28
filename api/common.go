package api

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"github.com/smallsixFight/hyakkei_blog/util/file_generator"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func VisitorAdd(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	ip := ctx.ClientIP()
	if ctx.GetHeader("token") != "" && !util.Cache.GetBool("visitor_"+ip) {
		util.Cache.Set("visitor_"+ip, true, time.Minute*5)
		util.Cache.Set(model.VisitorCount, util.Cache.GetInt32(model.VisitorCount)+1, 0)
	}
	reply.SetSuccess(true).SetMessage("ok")
}

func UploadAttachment(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	ft := ctx.PostForm("type")
	maxSize := getAttachSizeLimit(ft)
	if maxSize == 0 {
		reply.SetMessage("类型参数错误")
		return
	}
	file, err := ctx.FormFile("file")
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	if file.Size > maxSize {
		reply.SetMessage("文件过大，无法上传")
		return
	}
	saveDir := filepath.Join(service.GetSysConfig().SavePath, "hyakkei", getSavePath(ft))
	if !util.FileIsExist(saveDir) {
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			reply.SetMessage("存储目录创建失败: " + err.Error())
			return
		}
	}
	if ft == "favicon" && file.Filename != "favicon.ico" {
		reply.SetMessage("favicon 文件名必须是 favicon.ico")
		return
	}
	if err := ctx.SaveUploadedFile(file, filepath.Join(saveDir, file.Filename)); err != nil {
		reply.SetMessage("文件保存失败: " + err.Error())
		return
	}
	// 如果是 logo，更新文件信息，重新生成 header 文件
	cfg := service.GetSysConfig()
	if ft == "logo" && cfg.LogoName != file.Filename {
		cfg.LogoName = file.Filename
		go saveFilenameAndGenerateHeader(&cfg)
	}
	reply.SetSuccess(true).SetMessage("上传成功")
}

func saveFilenameAndGenerateHeader(cfg *model.SysSetting) {
	if err := service.SaveSysConfig(cfg); err != nil {
		logger.Warn("更新附件配置失败，", err.Error())
	} else if err = file_generator.GenerateHeaderFile(); err != nil {
		logger.Warn("重新生成 header 文件失败，", err.Error())
	}

}

func getSavePath(typ string) string {
	switch typ {
	case "logo", "favicon":
		return "assets/img"
	default:
		return "assets"
	}
}

func getAttachSizeLimit(typ string) int64 {
	const m = 1024 * 1024
	switch typ {
	case "logo":
		return 5 * m
	case "favicon":
		return 512 * 1024
	default:
		return 0
	}
}

// 预览 markdown 生成的 html 格式
func PreviewMarkdown(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	content := struct {
		Data string `json:"data"`
	}{}
	if err := ctx.BindJSON(&content); err != nil {
		reply.SetMessage("提交参数错误")
		return
	}
	var b bytes.Buffer
	if err := util.GetMDHandle().Convert([]byte(content.Data), &b); err != nil {
		reply.SetMessage("文本转换失败，" + err.Error())
		return
	}
	reply.SetData(b.String()).SetSuccess(true).SetMessage("ok")
}

func handlePostParam(post *model.Post, typ string) error {
	post.Title = strings.TrimSpace(post.Title)
	post.Slug = strings.TrimSpace(post.Slug)
	post.Typ = strings.TrimSpace(post.Typ)
	if post.Title == "" {
		return errors.New("标题不能为空")
	} else if post.Status != model.Publish && post.Status != model.Draft {
		return errors.New("状态错误")
	} else if post.Typ != typ {
		return errors.New("报文类型错误")
	}
	return nil
}

func getPostName(article *model.Post) string {
	if article.Slug != "" {
		return article.Slug
	}
	return article.Title
}
