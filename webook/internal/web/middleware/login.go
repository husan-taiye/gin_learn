package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//if ctx.Request.URL.Path == "/user/login" ||
		//	ctx.Request.URL.Path == "/user/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		//if sess==nil {
		//	// 没有登录
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now().UnixMilli()
		updateTime := sess.Get("update_time")
		sess.Options(sessions.Options{
			MaxAge: 20,
		})
		//sess.Save()
		//sess.Set("userId", id)
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		// updateTime 是有的
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		if now-updateTimeVal > 10*1000 {
			sess.Set("update_time", now)
			sess.Save()
		}
	}
}
