package user

import (
	"bytes"
	"context"
	"errors"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/service"
	svcmocks "gin_learn/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return userSvc
			},
			reqBody: `
{
    "email": "cb918551@qq.com",
    "password": "Helloworld123",
    "confirmPassword": "Helloworld123"
}
`,
			wantCode: 200,
			wantBody: "注册成功",
		},
		{
			name: "参数不对， bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return userSvc
			},
			reqBody: `
{
    "email": "cb918551@qq.com",
    "password": "Helloworld123",
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
{
    "email": "cb918551@q",
    "password": "Helloworld123",
    "confirmPassword": "Helloworld123"
}
`,
			wantCode: 200,
			wantBody: "邮箱格式错误",
		},
		{
			name: "邮箱重复",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)
				return userSvc
			},
			reqBody: `
{
    "email": "fasfas@qq.com",
    "password": "Helloworld123",
    "confirmPassword": "Helloworld123"
}
`,
			wantCode: 200,
			wantBody: "邮箱重复",
		},
	}
	//resp.Body
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterUserRouter(server)
			req, _err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, _err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}

func TestMock(t *testing.T) {
	// 初始化控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 创建模拟对象
	usersvc := svcmocks.NewMockUserService(ctrl)
	// 设计模拟调用
	// 先调用expect；	 调用同名方法，传入模拟条件；		指定返回值
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	err := usersvc.SignUp(context.Background(), domain.User{Email: "123@qq.com"})
	t.Log(err)
}
