package integration

import (
	"bytes"
	"encoding/json"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository/dao/article"
	"gin_learn/webook/internal/web/integration/startup"
	ijwt "gin_learn/webook/internal/web/jwt"
	"gin_learn/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	// 初始化全部内容
	//s.server = startup.InitWebServer()
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	s.db = ioc.InitDB(ioc.InitLogger())
	//artHdl := article.NewArticleHandler(service.NewArticleService(), ioc.InitLogger())
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TearDownTest() {
	// 清空articles
	s.db.Exec("TRUNCATE TABLE articles")
}

func (s *ArticleTestSuite) TestABC() {
	s.T().Log("这是测试套件")
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string

		// 集成测试准备数据
		before func(t *testing.T)
		// 集成测试验证数据
		after func(t *testing.T)
		// 预期输入
		art Article

		// HTTP 响应码
		wantCode int
		// HTTP响应带上 id
		wantRes Result[int64]
	}{
		{
			name: "新建帖子-保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				// 验证数据库
				var art article.Article
				err := s.db.Where("title=?", "my first title").First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Utime = 0
				art.Ctime = 0
				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "my first title",
					Content:  "xxxxx",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, art)
			},
			art: Article{
				Title:   "my first title",
				Content: "xxxxx",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data:    1,
				Msg:     "OK",
				Success: true,
			},
		},
		{
			name: "修改已有帖子-并保存",
			before: func(t *testing.T) {
				err := s.db.Create(article.Article{
					Id:       2,
					Title:    "已有帖子",
					Content:  "帖子内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
					Ctime:    1000,
					Utime:    1000,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art article.Article
				err := s.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 1000)
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					Ctime:    1000,
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data:    2,
				Msg:     "OK",
				Success: true,
			},
		},
		{
			name: "修改别人的帖子-并保存",
			before: func(t *testing.T) {
				err := s.db.Create(article.Article{
					Id:       3,
					Title:    "已有帖子",
					Content:  "帖子内容",
					AuthorId: 789,
					Ctime:    1000,
					Utime:    1000,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art article.Article
				err := s.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)

				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "已有帖子",
					Content:  "帖子内容",
					Ctime:    1000,
					Utime:    1000,
					AuthorId: 789,
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:     "系统错误",
				Code:    5,
				Success: false,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// 构造请求
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/article/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			//req.Header.Set("Authorization", "application/json")
			// 这里你就可以继续使用 req

			// 执行
			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			s.server.ServeHTTP(resp, req)

			// 验证数据
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)
		})
	}
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
	Success bool   `json:"success"`
}
