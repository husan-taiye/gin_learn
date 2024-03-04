package article

import (
	"bytes"
	"encoding/json"
	"errors"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/service"
	svcmocks "gin_learn/webook/internal/service/mocks"
	ijwt "gin_learn/webook/internal/web/jwt"
	"gin_learn/webook/internal/web/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {

	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  utils.Result
	}{
		{
			name: "新建并发表帖子",
			reqBody: `
{
	"title": "我的标题",
"content": "我的内容"
}`,
			wantCode: 200,
			wantRes: utils.Result{
				Msg:     "OK",
				Success: true,
				Data:    float64(1),
			},
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
		},
		{
			name: "publish发表失败",
			reqBody: `
{
	"title": "我的标题",
"content": "我的内容"
}`,
			wantCode: 200,
			wantRes: utils.Result{
				Code:    5,
				Msg:     "系统错误",
				Success: false,
			},
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("发表失败"))
				return svc
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", &ijwt.UserClaims{
					Uid: 123,
				})
			})
			h := NewArticleHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)
			req, _err := http.NewRequest(http.MethodPost, "/article/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, _err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes utils.Result
			err := json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)

			//assert.Equal(t, tc.wantBody, resp.Body)

		})
	}
}
