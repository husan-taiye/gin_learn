package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	//db := initDB()
	//redisClient := initCache()
	//userHandler := initUser(db, redisClient)
	//server := web.InitWebserver()
	//
	//ug := web.DispatchRoutes(server)
	//userHandler.RegisterUserRouter(ug)
	err := os.Setenv("WECHAT_APP_ID", "27017")
	err = os.Setenv("WECHAT_APP_SECRET", "27017")
	server := InitWebServer()
	//server := gin.Default()
	server.GET("hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello~")
	})
	err = server.Run(":8080")
	if err != nil {
		return
	}
}

//func initUser(db *gorm.DB, client *redis.Client) *user.UserHandler {
//	ud := dao2.NewUserDao(db)
//	userCache := cache.NewUserCache(client)
//	repo := repository.NewUserRepository(ud, userCache)
//	svc := wechat.NewUserService(repo)
//	// code初始化
//	codeCache := cache.NewCodeCache(client)
//	codeRepo := repository.NewCodeRepository(codeCache)
//	smsSvc := memory.NewService()
//	codeSvc := wechat.NewCodeService(codeRepo, smsSvc, "22321")
//	u := user.NewUserHandler(svc, codeSvc)
//	return u
//}

//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
//	if err != nil {
//		// 只在初始化过程中panic
//		// panic相当于整个goroutine结束
//		// 一旦初始化过程出错，就不再继续
//		panic(err)
//	}
//	err = dao2.InitTables(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

//func initCache() redis.Cmdable {
//	return redis.NewClient(&redis.Options{
//		Addr:     config.Config.Redis.Addr,
//		Password: "", // 密码
//		DB:       0,  // 数据库
//		PoolSize: 20, // 连接池大小
//	})
//}
