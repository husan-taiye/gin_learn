package repository

import (
	"context"
	"gin_learn/webook/domain"
	"gin_learn/webook/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (up *UserRepository) Create(ctx context.Context, u domain.User) error {
	return up.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 在这里操作缓存
}

func (up *UserRepository) FindById(int64) {

}
