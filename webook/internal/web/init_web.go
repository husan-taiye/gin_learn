package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	DispatchRoutes(server)
	return server
}
