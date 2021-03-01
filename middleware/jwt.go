package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"net/http"
	"strings"
	"time"
)

func JWTVerify(c *gin.Context) {
	reply := model.Reply{Code: -1, Message: "授权验证失败"}
	token := strings.TrimSpace(c.GetHeader("accessToken"))
	if token == "" || len(token) < 32 {
		c.JSON(http.StatusUnauthorized, &reply)
		c.Abort()
		return
	}
	jwtToken, err := parseToken(token)
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

type MyCustomClaims struct {
	jwt.StandardClaims
	CustomParams map[string]interface{} `json:"custom_params"`
}

func CreateToken(customClaims map[string]interface{}, expired time.Duration) (token string, err error) {
	claims := MyCustomClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    "hyakkeiBlog",
			Subject:   "admin",
			Audience:  "",
			ExpiresAt: time.Now().Add(expired).Unix(),
		},
		customClaims,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(service.GetSysConfig().TokenSecret))
}

func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (key interface{}, err error) {
		return []byte(service.GetSysConfig().TokenSecret), nil
	})
}
