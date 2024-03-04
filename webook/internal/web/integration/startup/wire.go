//go:build wireinject

package startup

import (
	"gin_learn/webook/internal/repository"
	article2 "gin_learn/webook/internal/repository/article"
	"gin_learn/webook/internal/repository/cache"
	"gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web"
	"gin_learn/webook/internal/web/article"
	ijwt "gin_learn/webook/internal/web/jwt"
	"gin_learn/webook/internal/web/user"
	"gin_learn/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// 最基础的第三方依赖
var thirdProvider = wire.NewSet(ioc.InitDB, ioc.InitRedis, ioc.InitLogger)

// 用户服务提供
var userSvcProvider = wire.NewSet(
	dao.NewUserDao,
	cache.NewUserCache,
	repository.NewUserRepository,
	ioc.InitUserService,
	//service.NewUserService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		// dao,cache
		cache.NewCodeCache,
		dao.NewArticleDAO,
		// repo
		repository.NewCodeRepository,
		article2.NewArticleRepository,
		// svc
		service.NewCodeService,
		service.NewArticleService,
		//service.NewUserService,
		ioc.InitSMSService,
		// handler
		user.NewUserHandler,
		web.NewOAuth2WechatHandler,
		article.NewArticleHandler,
		ijwt.NewRedisJWTHandler,

		//ioc.InitUserService,
		ioc.InitGin,

		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,
		ioc.InitMiddlewares,
		ioc.InitTpl,
	)
	return new(gin.Engine)
}

func InitArticleHandler() *article.ArticleHandler {
	wire.Build(thirdProvider,
		article2.NewArticleRepository,
		dao.NewArticleDAO,
		service.NewArticleService,
		article.NewArticleHandler,
	)
	return &article.ArticleHandler{}
}
