package repository

import (
	"context"
	"gin_learn/webook/domain"
	"gin_learn/webook/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFount = dao.ErrUserNotFount

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (up *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := up.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
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
