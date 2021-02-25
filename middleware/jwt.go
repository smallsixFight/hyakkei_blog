package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/util"
	"net/http"
	"strings"
)

func JWTVerify(c *gin.Context) {
	reply := model.Reply{Code: -1, Message: "授权验证失败"}
	token := strings.TrimSpace(c.GetHeader("accessToken"))
	if token == "" || len(token) < 32 {
		c.JSON(http.StatusUnauthorized, &reply)
		c.Abort()
		return
	}
	jwtToken, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &reply)
		c.Abort()
		return
	}
	m, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		reply.SetMessage("获取授权信息失败")
		c.JSON(http.StatusUnauthorized, &reply)
		c.Abort()
	}
	customMap, ok := m["custom_params"].(map[string]interface{})
	if !ok {
		reply.SetMessage("获取授权信息失败")
		c.JSON(http.StatusUnauthorized, &reply)
		c.Abort()
	}
	c.Set("user_id", customMap["user_id"])
	c.Set("username", customMap["username"])
	c.Next()
}
