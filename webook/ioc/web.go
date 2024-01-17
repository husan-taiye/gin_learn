package ioc

import (
	"gin_learn/webook/internal/web"
	"gin_learn/webook/internal/web/middleware"
	"gin_learn/webook/internal/web/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitGin(middles []gin.HandlerFunc, handler *user.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middles...)
	ug := web.DispatchRoutes(server)
	handler.RegisterUserRouter(ug)
	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHandler(),
		middleware.NewloginJWTMiddlewareBuilder().
			IgnorePaths("/user/signup").
			IgnorePaths("/user/login_sms/code/send").
			IgnorePaths("/user/login_sms").
			IgnorePaths("/user/login").Build(),
	}
}

func InitTpl() string {
	return "321312"
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:  []string{"PUT", "PATCH", "OPTIONS", "POST"},
		AllowHeaders:  []string{"Origin", "Authorization"},
		ExposeHeaders: []string{"Authorization", "Content-Type", "X-jwt-token"},
		// 是否允许带cookie
		AllowCredentials: true,
		// 下面两个都可以，二选一
		// AllowOrigins:  []string{"http://localhost:3000"},
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "localhost")
		},
		MaxAge: 12 * time.Hour,
	})
}
