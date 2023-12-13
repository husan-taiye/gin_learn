package web

import (
	"gin_learn/webook/internal/web/user"
	//"gin_learn/webook/service"
	"github.com/gin-gonic/gin"
)

func DispatchRoutes(server *gin.Engine) {
	// user 路由
	ug := server.Group("/user")
	//svc = server.User
	u := user.NewUserHandler()
	u.RegisterUserRouter(ug)
}
