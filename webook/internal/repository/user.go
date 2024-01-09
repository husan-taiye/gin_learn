package repository

import (
	"context"
	"errors"
	"fmt"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository/cache"
	"gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/utils"
	"time"
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFount = dao.ErrUserNotFount

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
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
	cacheUp, err := ur.cache.Get(ctx, userId)
	//switch err {
	//case nil:
	//	return cacheUp, nil
	//case ErrUserNotFount:
	//	up, err := ur.dao.FindProfileByUserId(ctx, userId)
	//	if err != nil {
	//		return domain.UserProfile{}, err
	//	}
	//
	//	toCacheUp := domain.UserProfile{
	//		UserId:   up.UserId,
	//		Nickname: up.Nickname,
	//		Birthday: time.UnixMilli(up.Birthday).Format(utils.TimeTemplate3),
	//		Profile:  up.Profile,
	//	}
	//	err = ur.cache.Set(ctx, toCacheUp)
	//	if err != nil {
	//		// 打日志，做监控
	//		fmt.Sprintln("同步缓存失败")
	//	}
	//	return toCacheUp, err
	//default:
	//	return domain.UserProfile{}, nil
	//}
	if err == nil {
		return cacheUp, nil
	}
	if errors.Is(err, cache.ErrUserNotFound) {
		up, err := ur.dao.FindProfileByUserId(ctx, userId)
		if err != nil {
			return domain.UserProfile{}, err
		}

		toCacheUp := domain.UserProfile{
			UserId:   up.UserId,
			Nickname: up.Nickname,
			Birthday: time.UnixMilli(up.Birthday).Format(utils.TimeTemplate3),
			Profile:  up.Profile,
		}
		err = ur.cache.Set(ctx, toCacheUp)
		if err != nil {
			// 打日志，做监控
			fmt.Sprintln("同步缓存失败")
		}
		return toCacheUp, err
	}
	// redis 其他错误，可能是崩了
	// 保护数据库，数据库查询限流或者其他方式等等
	// 这里暂时就先不查询了
	return domain.UserProfile{}, nil
}
