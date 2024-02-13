package service

import (
	"context"
	"errors"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserDuplicate = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
var ErrUserNotFount = repository.ErrUserNotFount

type UserService interface {
	Login(ctx context.Context, user domain.User) (domain.User, error)
	SignUp(ctx context.Context, user domain.User) error
	Edit(ctx context.Context, up domain.UserProfile) error
	Profile(ctx context.Context, userId int64) (domain.UserProfile, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, openId string) (domain.User, error)
}

type RepoUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &RepoUserService{
		repo: repo,
	}
}

func (svc *RepoUserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	// 先找到用户
	findUser, err := svc.repo.FindByEmail(ctx, user.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return findUser, nil
}

func (svc *RepoUserService) SignUp(ctx context.Context, user domain.User) error {
	// 加密放这里
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//bcrypt.CompareHashAndPassword()
	if err != nil {
		return err
	}
	user.Password = string(hash)
	// 存储
	return svc.repo.Create(ctx, user)
}

func (svc *RepoUserService) Edit(ctx context.Context, up domain.UserProfile) error {
	return svc.repo.Update(ctx, up)
}

func (svc *RepoUserService) Profile(ctx context.Context, userId int64) (domain.UserProfile, error) {
	return svc.repo.FindProfileByUserId(ctx, userId)
}

func (svc *RepoUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFount) {
		return u, err
	}
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return u, err
	}
	// 可能会碰到主从延迟问题
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *RepoUserService) FindOrCreateByWechat(ctx context.Context, openId string) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, openId)
	if !errors.Is(err, repository.ErrUserNotFount) {
		return u, err
	}
	u = domain.User{
		OpenId: openId,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return u, err
	}
	// 可能会碰到主从延迟问题
	return svc.repo.FindByWechat(ctx, openId)
}

//type PathsDownGrade struct {
//	Quick
//}

func PathsDownGrade(ctx context.Context, quick, slow func()) {
	quick()
	if ctx.Value("降级") == "true" {
		return
	}
	slow()
}
