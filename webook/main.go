package main

import (
	"gin_learn/webook/internal/config"
	"gin_learn/webook/internal/repository"
	dao2 "gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web"
	"gin_learn/webook/internal/web/user"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	db := initDB()
	userHandler := initUser(db)
	server := web.InitWebserver()

	ug := web.DispatchRoutes(server)
	userHandler.RegisterUserRouter(ug)
	//server := gin.Default()
	server.GET("hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello~")
	})
	err := server.Run(":8081")
	if err != nil {
		return
	}
}

func initUser(db *gorm.DB) *user.UserHandler {
	ud := dao2.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := user.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 只在初始化过程中panic
		// panic相当于整个goroutine结束
		// 一旦初始化过程出错，就不再继续
		panic(err)
	}
	err = dao2.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
