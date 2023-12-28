package middleware

import (
	"encoding/gob"
	"gin_learn/webook/internal/web/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

// jwt 登录校验
type loginJWTMiddlewareBuilder struct {
	paths []string
}

func NewloginJWTMiddlewareBuilder() *loginJWTMiddlewareBuilder {
	return &loginJWTMiddlewareBuilder{}
}
func (l *loginJWTMiddlewareBuilder) IgnorePaths(path string) *loginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *loginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := &user.UserClaims{}
		// ParseWithClaims一定要穿指针
		token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("r4BKnmqBgWhnudRc4xufW9f97ODTqX10"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < 50*time.Second {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err := token.SignedString([]byte("r4BKnmqBgWhnudRc4xufW9f97ODTqX10"))
			if err != nil {
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)
	}
}
