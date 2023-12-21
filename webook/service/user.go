package service

import (
	"context"
	"errors"
	"gin_learn/webook/domain"
	"gin_learn/webook/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
var ErrUserNotFount = repository.ErrUserNotFount

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
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

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
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

func (svc *UserService) Edit(ctx context.Context, up domain.UserProfile) error {
	return svc.repo.Update(ctx, up)
}

func (svc *UserService) Profile(ctx context.Context, userId int64) (domain.UserProfile, error) {
	return svc.repo.FindProfileByUserId(ctx, userId)
}
