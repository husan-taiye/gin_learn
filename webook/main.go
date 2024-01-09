package main

import (
	"gin_learn/webook/internal/config"
	"gin_learn/webook/internal/repository"
	"gin_learn/webook/internal/repository/cache"
	dao2 "gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web"
	"gin_learn/webook/internal/web/user"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	db := initDB()
	redisClient := initCache()
	userHandler := initUser(db, redisClient)
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

func initUser(db *gorm.DB, client *redis.Client) *user.UserHandler {
	ud := dao2.NewUserDao(db)
	userCache := cache.NewUserCache(client)
	repo := repository.NewUserRepository(ud, userCache)
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

func initCache() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: "", // 密码
		DB:       0,  // 数据库
		PoolSize: 20, // 连接池大小
	})
}
