package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository/cache"
	"gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/utils"
	"time"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	Update(ctx context.Context, up domain.UserProfile) error
	FindProfileByUserId(ctx context.Context, userId int64) (domain.UserProfile, error)
}

type DCUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFount = dao.ErrUserNotFount

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &DCUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (ur *DCUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.modelToDomain(u), nil
}

func (ur *DCUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.modelToDomain(u), nil
}

func (ur *DCUserRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	u, err := ur.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return ur.modelToDomain(u), nil
}

func (ur *DCUserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domainToModel(u))
	// 在这里操作缓存
}

func (ur *DCUserRepository) Update(ctx context.Context, up domain.UserProfile) error {
	BirthdayStamp, _ := time.ParseInLocation(utils.TimeTemplate3, up.Birthday, time.Local)
	return ur.dao.Update(ctx, dao.UserProfile{
		UserId:   up.UserId,
		Nickname: up.Nickname,
		Profile:  up.Profile,
		Birthday: BirthdayStamp.UnixMilli(),
	})
}

func (ur *DCUserRepository) FindProfileByUserId(ctx context.Context, userId int64) (domain.UserProfile, error) {
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
			//Birthday: time.UnixMilli(up.Birthday).Format(time.DateOnly),
			Profile: up.Profile,
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

func (ur *DCUserRepository) domainToModel(du domain.User) dao.User {
	return dao.User{
		Id: du.Id,
		Email: sql.NullString{
			String: du.Email,
			Valid:  du.Email != "",
		},
		Phone: sql.NullString{
			String: du.Phone,
			Valid:  du.Phone != "",
		},
		OpenId: sql.NullString{
			String: du.OpenId,
			Valid:  du.OpenId != "",
		},
		Password: du.Password,
	}
}

func (ur *DCUserRepository) modelToDomain(ud dao.User) domain.User {
	return domain.User{
		Id:       ud.Id,
		Email:    ud.Email.String,
		Phone:    ud.Phone.String,
		Password: ud.Password,
		OpenId:   ud.OpenId.String,
	}
}
