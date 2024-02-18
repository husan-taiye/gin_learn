// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gin_learn/webook/internal/repository"
	"gin_learn/webook/internal/repository/cache"
	"gin_learn/webook/internal/repository/dao"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web"
	"gin_learn/webook/internal/web/jwt"
	"gin_learn/webook/internal/web/user"
	"gin_learn/webook/ioc"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitMiddlewares(handler)
	db := ioc.InitDB()
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	string2 := ioc.InitTpl()
	codeService := service.NewCodeService(codeRepository, smsService, string2)
	userHandler := user.NewUserHandler(userService, codeService, handler)
	wechatService := ioc.InitWechatService()
	wechatHandlerConfig := ioc.NewWechatHandlerConfig()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, wechatHandlerConfig, handler)
	engine := ioc.InitGin(v, userHandler, oAuth2WechatHandler)
	return engine
}
