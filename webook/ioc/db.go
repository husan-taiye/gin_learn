package ioc

import (
	dao2 "gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func InitDB(logger logger.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:3309)/webook",
	}
	err := viper.UnmarshalKey("db.mysql", &cfg)
	//dsn := viper.GetString("db.mysql.dsn")
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		// 缺了一个writer
		Logger: glogger.New(gormLoggerFunc(logger.Debug), glogger.Config{
			// 慢查询阈值，只有执行时间超过了这个阈值才会使用
			// 50ms,100ms
			// sql查询必须命中索引，最好就是走一次磁盘IO
			// 一次磁盘 IO 是不到 10ms
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  glogger.Info,
		}),
	})
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

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
