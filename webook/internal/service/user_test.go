package service

import (
	"context"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository"
	repomocks "gin_learn/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

func TestRepoUserService_Login(t *testing.T) {
	// 测试用公共时间
	//now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		//ctx      context.Context
		user     domain.User
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$AKZyoMYcV9aSFAByysGUpOlGAcxG45MFhHg0XL7rvv2BJxN3zYHaK",
						Phone:    "13953859925",
					}, nil)
				return repo
			},
			user: domain.User{
				Email:    "123@qq.com",
				Password: "l913687515",
			},
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$AKZyoMYcV9aSFAByysGUpOlGAcxG45MFhHg0XL7rvv2BJxN3zYHaK",
				Phone:    "13953859925",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, gorm.ErrRecordNotFound)
				return repo
			},
			user: domain.User{
				Email:    "123@qq.com",
				Password: "l913687515",
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			user, err := svc.Login(context.Background(), tc.user)

			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}

func TestEncrypted(t *testing.T) {
	res, _ := bcrypt.GenerateFromPassword([]byte("l913687515"), bcrypt.DefaultCost)
	t.Log(string(res))
}
