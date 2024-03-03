package article

import (
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web/handler"
	ijwt "gin_learn/webook/internal/web/jwt"
	"gin_learn/webook/internal/web/utils"
	"gin_learn/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ handler.Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.Logger
}

func NewArticleHandler(svc service.ArticleService, logger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: logger,
	}
}

func (art *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/article")
	g.POST("/edit", art.Edit)
}

func (art *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content""`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 获取用户
	c, _ := ctx.Get("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 监控输出
		ctx.JSON(http.StatusOK, utils.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		art.logger.Error("未发现用户的 session 信息")
		return
	}
	userId := claims.Uid

	// 检测输入
	// 调用service
	id, err := art.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: userId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志
		art.logger.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Msg:     "OK",
		Data:    id,
		Success: true,
	})
}
