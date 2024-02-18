package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uId int64) error
	SetJWTToken(ctx *gin.Context, uId int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 生命要放进token里面的数据
	Uid       int64
	Ssid      string
	UserAgent string
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}
