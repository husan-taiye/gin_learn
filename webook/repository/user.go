package repository

import (
	"context"
	"gin_learn/webook/domain"
	"gin_learn/webook/repository/dao"
	"gin_learn/webook/utils"
	"time"
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

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 在这里操作缓存
}

func (ur *UserRepository) Update(ctx context.Context, up domain.UserProfile) error {
	BirthdayStamp, _ := time.ParseInLocation(utils.TimeTemplate3, up.Birthday, time.Local)
	return ur.dao.Update(ctx, dao.UserProfile{
		UserId:   up.UserId,
		Nickname: up.Nickname,
		Profile:  up.Profile,
		Birthday: BirthdayStamp.UnixMilli(),
	})
}

func (ur *UserRepository) FindProfileByUserId(ctx context.Context, userId int64) (domain.UserProfile, error) {
	up, err := ur.dao.FindProfileByUserId(ctx, userId)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return domain.UserProfile{
		UserId:   up.UserId,
		Nickname: up.Nickname,
		Birthday: time.UnixMilli(up.Birthday).Format(utils.TimeTemplate3),
		Profile:  up.Profile,
	}, nil
}
