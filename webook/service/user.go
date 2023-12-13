package service

import (
	"context"
	"gin_learn/webook/domain"
)

type UserService struct {
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	return nil
}
