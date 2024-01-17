package repository

import "gin_learn/wire/repository/dao"

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(ud *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: ud,
	}
}
