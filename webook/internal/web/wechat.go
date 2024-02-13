package web

import (
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/service/oauth2/wechat"
	"gin_learn/webook/internal/web/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	utils.JwtHandler
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/oauth2")
	g.GET("/auth_url", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	// 存储state
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Result{
			Code: 5,
			Msg:  "构造扫码登录url失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{Data: url})
	return
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 验证微信的code
	code := ctx.Query("code")
	state := ctx.Query("state")
	// 校验state
	info, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 5, Msg: "系统错误"})
		return
	}
	// 从userSvc里面拿uid
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info.OpenId)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 5, Msg: "系统错误"})
		return
	}

	err = h.SetJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{Msg: "OK"})
	return
}

//type OAuth2Handler struct {
//}
//
//func (h *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
//	// 统一处理所有的oauth2
//	g := server.Group("/oauth2")
//	g.GET("/:platform/auth_url", h.AuthURL)
//	g.Any("/:platform/callback", h.Callback)
//}
//
//func (h *OAuth2Handler) AuthURL(ctx *gin.Context) {
//	platform := ctx.Param("platform")
//	switch platform {
//	case "oauth2":
//		h.wechatService.AuthURL
//	}
//}
//
//func (h *OAuth2Handler) Callback(ctx *gin.Context) {
//	return
//}
