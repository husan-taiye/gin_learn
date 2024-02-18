package web

import (
	"errors"
	"fmt"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/service/oauth2/wechat"
	"gin_learn/webook/internal/web/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	utils.JwtHandler
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("r2BKnmqBgWhnudRc4xufW9f97ODTqX12"),
		cfg:      cfg,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/auth_url", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	// 存储state
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Result{
			Code: 5,
			Msg:  "构造扫码登录url失败",
		})
		return
	}
	if err = h.setStateCookie(ctx, state, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{Data: url})
	return
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string, err error) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		state: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr, 600,
		"oauth2/wechat/callback", "", h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 验证微信的code
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 5, Msg: "系统错误"})
		return
	}
	// 校验code
	info, err := h.svc.VerifyCode(ctx, code)
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
	err = h.SetRefreshToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{Msg: "OK"})
	return
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	// 校验state
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		// 做好监控
		// 有人乱搞
		return fmt.Errorf("拿不到state的cookie， %w", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		//ctx.JSON(http.StatusOK, utils.Result{Code: 4, Msg: "登录失败"})
		return fmt.Errorf("解析token失败， %w", err)
	}
	if sc.state != state {
		//ctx.JSON(http.StatusOK, utils.Result{Code: 4, Msg: "登录失败"})
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaims struct {
	state string
	jwt.RegisteredClaims
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
