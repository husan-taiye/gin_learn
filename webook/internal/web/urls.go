package web

import (
	"github.com/gin-gonic/gin"
)

func DispatchRoutes(server *gin.Engine) *gin.RouterGroup {
	// user 路由
	ug := server.Group("/user")
	//svc = server.User
	return ug
}
