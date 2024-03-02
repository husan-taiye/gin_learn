package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
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
	//initViper()
	//initViperV1()
	//initViperReader()
	initViperV3Remote()
	initLogger()
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

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	// 实时监听配置变更
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func initViperV1() {
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:3306)/webook"
redis:
  addr: "localhost:30003"
`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initViperV2() {
	cflle := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cflle)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// initViperV3Remote viper远程配置中心 etcd
func initViperV3Remote() {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	// 监听远程变更配置
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}

	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("replace之前")
	// 如果不replace 什么都打不出来
	zap.ReplaceGlobals(logger)
	zap.L().Info("replace之后")
	// zap第一个用法
	zap.L().Error("发送验证码失败", zap.Error(err))
	zap.L().Info("",
		zap.Error(errors.New("日志初始化失败")),
		zap.String("phone", "123"),
	)
}
