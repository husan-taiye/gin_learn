package web

import (
	"gin_learn/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitWebserver() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{

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
	}))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("f4TgXINAWeleaJ3f70AI7J3vTKQtJjnO"), []byte("DYs8aItQBkFa9pw8KpK0AkRn7XPsPN1g"))
	if err != nil {
		panic(err)
	}
	//store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: 20,
	})
	server.Use(sessions.Sessions("mysession", store))

	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/user/signup").
	//	IgnorePaths("/user/login").Build())
	server.Use(middleware.NewloginJWTMiddlewareBuilder().
		IgnorePaths("/user/signup").
		IgnorePaths("/user/login").Build())

	return server
}
