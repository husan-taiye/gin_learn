package ioc

import (
	"gin_learn/webook/internal/repository"
	"gin_learn/webook/internal/service"
	"go.uber.org/zap"
)

func InitUserService(repo repository.UserRepository) service.UserService {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return service.NewUserService(repo, l)
}
