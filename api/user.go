package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/middleware"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"net/http"
	"strings"
	"time"
)

type loginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	param := loginParam{}
	key := fmt.Sprintf("login_fail_%s", ctx.ClientIP())
	failCount := util.Cache.GetInt32(key)
	if failCount >= 50 {
		reply.SetMessage("登录失败多次，请稍后再尝试登录")
		return
	}
	if err := ctx.BindJSON(&param); err != nil {
		logger.Println("请求参数错误，" + err.Error())
		reply.Message = "请求参数错误"
		return
	}
	if strings.TrimSpace(param.Username) == "" || strings.TrimSpace(param.Password) == "" {
		reply.SetMessage("用户名/密码不能为空")
		return
	}
	token, err := VerifyLoginInfo(&param)
	if err != nil {
		logger.Warnf("登录失败，用户名：%s，密码：%s", param.Username, param.Password)
		util.Cache.Set(key, failCount+1, time.Minute*30)
		reply.SetMessage(err.Error())
		return
	}

	reply.SetSuccess(true).SetToken(token).SetMessage("登录成功")
}

func VerifyLoginInfo(param *loginParam) (token string, err error) {
	cfg := service.GetSysConfig()
	// 密码加密验证
	result, err := util.MD5Encrypt(param.Password, cfg.Salt)
	if err != nil {
		return
	}
	if result != cfg.Password {
		return "", errors.New("帐号/密码错误")
	}
	// 生成 token
	claims := make(map[string]interface{})
	claims["username"] = param.Username
	token, err = middleware.CreateToken(claims, time.Hour*6)
	if err != nil {
		return
	}
	return
}
