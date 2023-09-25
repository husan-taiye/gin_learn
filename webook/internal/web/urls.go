package web

import (
	"gin_learn/webook/internal/web/user"
	"github.com/gin-gonic/gin"
)

func DispatchRoutes(server *gin.Engine) {
	// user 路由
	ug := server.Group("/user")
	u := user.NewUserHandler()
	u.RegisterUserRouter(ug)
}
