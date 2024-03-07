package article

import (
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web/handler"
	ijwt "gin_learn/webook/internal/web/jwt"
	"gin_learn/webook/internal/web/utils"
	"gin_learn/webook/pkg/ginx"
	"gin_learn/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	g.POST("/publish", art.Publish)
	g.POST("/withdraw", art.Withdraw)
	// get获取
	g.GET("/list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](art.List))
}

func (art *ArticleHandler) Edit(ctx *gin.Context) {

	var req ArticleReq
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

	// 检测输入
	// 调用service
	id, err := art.svc.Save(ctx, req.toDomain(claims.Uid))
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

func (art *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
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

	id, err := art.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code:    5,
			Msg:     "系统错误",
			Success: false,
		})
		// 打日志
		art.logger.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Msg:     "OK",
		Data:    id,
		Success: true,
	})
}

func (art *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
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

	err := art.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code:    5,
			Msg:     "系统错误",
			Success: false,
		})
		// 打日志
		art.logger.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{
		Msg:     "OK",
		Success: true,
	})

}

func (art *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (utils.Result, error) {
	res, err := art.svc.List(ctx, uc.Uid, req.Offset, req.Limit)
	if err != nil {
		return utils.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	// 列表页，不显示全文，只显示一个摘要
	// 前几句话
	return utils.Result{
		Data: slice.Map[domain.Article, ArticleVO](res,
			func(idx int, src domain.Article) ArticleVO {
				return ArticleVO{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Status:   src.Status.ToUint8(),
					Ctime:    src.Ctime.Format(time.DateTime),
					Utime:    src.Utime.Format(time.DateTime),
				}
			}),
	}, nil
}
