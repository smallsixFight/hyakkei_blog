package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

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
	return t.SignedString([]byte(""))
}

func ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (key interface{}, err error) {
		return []byte(""), nil
	})
}
