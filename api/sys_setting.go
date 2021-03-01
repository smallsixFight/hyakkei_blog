package api

import (
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util/file_generator"
	"net/http"
)

func GetSysSettingInfo(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	cfg := service.GetSysConfig()
	cfg.Password = ""
	cfg.Salt = ""
	cfg.TokenSecret = ""
	cfg.SavePath = ""
	reply.SetSuccess(true).SetData(&cfg).SetMessage("ok")
}

func UpdateSysSettingInfo(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	cfg := model.BaseSysSetting{}
	if err := ctx.BindJSON(&cfg); err != nil {
		logger.Println("提交参数错误: ", err.Error())
		reply.SetMessage("提交参数错误: " + err.Error())
		return
	}
	if err := service.SaveSysConfig(&model.SysSetting{BaseSysSetting: cfg}); err != nil {
		reply.SetMessage(err.Error())
		return
	}
	if err := file_generator.GenerateHeaderFile(); err != nil {
		reply.SetMessage(err.Error())
		return
	}
	reply.SetSuccess(true).SetMessage("更新成功")
}
