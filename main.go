package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.GET("/hello", func(c *gin.Context) { //静态路由
		oid := c.Query("id")
		c.String(http.StatusOK, "hello, gin  "+oid)
		//c.JSON(200, gin.H{
		//	"message": "pong",
		//})
	})
	server.POST("/hello/:name", func(c *gin.Context) { //参数路由
		name := c.Param("name")
		c.String(http.StatusOK, "hello, gin post"+name)

	})
	server.GET("/hello/*f", func(c *gin.Context) { //通配符路由
		page := c.Param(".f")

		c.String(http.StatusOK, "hello, 通配符 "+page)

	})
	err := server.Run(":8080")
	if err != nil {
		return
	} // 监听并在 0.0.0.0:8080 上启动服务
}
