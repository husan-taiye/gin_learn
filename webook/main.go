package main

import (
	"gin_learn/webook/internal/web"
	"gin_learn/webook/internal/web/user"
	"gin_learn/webook/repository"
	"gin_learn/webook/repository/dao"
	"gin_learn/webook/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"))
	if err != nil {
		// 只在初始化过程中panic
		// panic相当于整个goroutine结束
		// 一旦初始化过程出错，就不再继续
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := user.NewUserHandler(svc)

	server := web.RegisterRoutes(u)
	err = server.Run(":8000")
	if err != nil {
		return
	}
}
