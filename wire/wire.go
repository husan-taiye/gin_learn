//go:build wireinject

// 让wire来注入这里的代码

package wire

import (
	"gin_learn/wire/repository"
	"gin_learn/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}
