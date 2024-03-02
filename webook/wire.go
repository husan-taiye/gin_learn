//go:build wireinject

package main

import (
	"gin_learn/webook/internal/repository"
	"gin_learn/webook/internal/repository/cache"
	"gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web"
	ijwt "gin_learn/webook/internal/web/jwt"
	"gin_learn/webook/internal/web/user"
	"gin_learn/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的依赖
		ioc.InitDB, ioc.InitRedis,
		// dao,cache
		dao.NewUserDao,
		cache.NewCodeCache,
		cache.NewUserCache,
		// repo
		repository.NewUserRepository,
		repository.NewCodeRepository,
		// svc
		service.NewCodeService,
		//service.NewUserService,
		ioc.InitSMSService,
		// handler
		user.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitUserService,
		ioc.InitGin,
		ioc.InitLogger,
		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,
		ioc.InitMiddlewares,
		ioc.InitTpl,
	)
	return new(gin.Engine)
}
