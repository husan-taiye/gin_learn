package ioc

import (
	dao2 "gin_learn/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:3309)/webook",
	}
	err := viper.UnmarshalKey("db.mysql", &cfg)
	//dsn := viper.GetString("db.mysql.dsn")
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	db, err := gorm.Open(mysql.Open(cfg.DSN))
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
