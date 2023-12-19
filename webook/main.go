package main

import (
	"gin_learn/webook/internal/web"
	"gin_learn/webook/internal/web/user"
	"gin_learn/webook/repository"
	dao2 "gin_learn/webook/repository/dao"
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
	dao := dao2.NewUserDao(db)
	repo := repository.NewUserRepository(dao)
	svc := service.NewUserService(repo)
	u := user.NewUserHandler(svc)

	server := web.RegisterRoutes(u)
	err = server.Run(":8000")
	if err != nil {
		return
	}
}
