package ioc

import (
	"gin_learn/webook/internal/config"
	dao2 "gin_learn/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
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
