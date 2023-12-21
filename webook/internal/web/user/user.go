package user

import (
	"errors"
	"fmt"
	"gin_learn/webook/domain"
	"gin_learn/webook/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 定义所有有关user的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (user *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string
		ConfirmPassword string
		Password        string
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析数据到req里面
	// 解析错误会直接返回400的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//u := NewUserHandler()
	// email校验
	ok, err := user.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	// password校验
	ok, err = user.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式错误")
		return
	}
	// password与confirmPassword校验
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	// service 层
	err = user.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		//ctx.String(http.StatusOK, "邮箱重复")
		ctx.JSON(http.StatusOK, map[string]any{"msg": "邮箱重复", "success": false})
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	fmt.Printf("%v", req)
	ctx.String(http.StatusOK, "注册成功")
	return
}
func (user *UserHandler) Login(ctx *gin.Context) {
	// 定义接受结构体
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// 初始化接受结构体并赋值
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 往service层传值
	findUser, err := user.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "账号/邮箱或密码不对", "success": false})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{"msg": "系统错误", "success": false})
	}
	// 登录成功
	// 设置session
	sess := sessions.Default(ctx)
	// 需要放在session里面的值
	sess.Set("userId", findUser.Id)
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"msg": "登录成功", "success": true})
	return
}

func (user *UserHandler) Edit(ctx *gin.Context) {

}
func (user *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	ctx.JSON(http.StatusOK, map[string]any{"msg": "这是你的profile", "success": true, "userId": sess.Get("userId")})
}
