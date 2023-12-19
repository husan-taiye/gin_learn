package service

import (
	"context"
	"gin_learn/webook/domain"
	"gin_learn/webook/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	// 加密放哪里
	// 存储
	return svc.repo.Create(ctx, user)
}
