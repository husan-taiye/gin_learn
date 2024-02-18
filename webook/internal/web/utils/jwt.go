package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtHandler struct {
	// access token
	AtKey []byte
	// refresh token
	RtKey []byte
}

func NewJwtHandler() JwtHandler {
	return JwtHandler{
		RtKey: []byte("r4BKnmqBgWhnudRc4xufW9f97ODTqX10"),
		AtKey: []byte("r4BKnmqBgWhnudRc4xugd9f97ODTqX10"),
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 生命要放进token里面的数据
	Uid       int64
	UserAgent string
}

type RefreshClaims struct {
	Uid int64
	jwt.RegisteredClaims
}

func (h JwtHandler) SetJWTToken(ctx *gin.Context, uId int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		UserAgent: ctx.Request.UserAgent(),
		Uid:       uId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(h.AtKey)
	if err != nil {

		return err
	}
	fmt.Println(tokenStr)
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h JwtHandler) SetRefreshToken(ctx *gin.Context, uId int64) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		Uid: uId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(h.RtKey)
	if err != nil {

		return err
	}
	fmt.Println(tokenStr)
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}
