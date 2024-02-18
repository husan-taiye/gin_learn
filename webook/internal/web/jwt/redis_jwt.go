package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	RtKey = []byte("r4BKnmqBgWhnudRc4xufW9f97ODTqX10")
	AtKey = []byte("r4BKnmqBgWhnudRc4xugd9f97ODTqX10")
)

type RedisJWTHandler struct {
	//// access token
	//AtKey []byte
	//// refresh token
	//RtKey []byte
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (r RedisJWTHandler) SetLoginToken(ctx *gin.Context, uId int64) error {
	ssid := uuid.New().String()

	err := r.SetJWTToken(ctx, uId, ssid)
	if err != nil {
		return err
	}
	err = r.SetRefreshToken(ctx, uId, ssid)
	return err
}

func (r RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uId int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Ssid: ssid,
		Uid:  uId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {

		return err
	}
	fmt.Println(tokenStr)
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (r RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		return errors.New("解析claims失败 ")
	}
	return r.cmd.Set(ctx, fmt.Sprintf("user:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}

func (r RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := r.cmd.Exists(ctx, fmt.Sprintf("user:ssid:%s", ssid)).Result()
	return err
}

func (r RedisJWTHandler) SetJWTToken(ctx *gin.Context, uId int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
		Uid:       uId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {

		return err
	}
	fmt.Println(tokenStr)
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
