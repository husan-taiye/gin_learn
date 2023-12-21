package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitWebserver() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{

		AllowMethods:  []string{"PUT", "PATCH", "OPTIONS", "POST"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Authorization", "Content-Type"},
		// 是否允许带cookie
		AllowCredentials: true,
		// 下面两个都可以，二选一
		// AllowOrigins:  []string{"http://localhost:3000"},
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	return server
}
