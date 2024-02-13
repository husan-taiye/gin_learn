package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtHandler struct {
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 生命要放进token里面的数据
	Uid       int64
	UserAgent string
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
	tokenStr, err := token.SignedString([]byte("r4BKnmqBgWhnudRc4xufW9f97ODTqX10"))
	if err != nil {

		return err
	}
	fmt.Println(tokenStr)
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
