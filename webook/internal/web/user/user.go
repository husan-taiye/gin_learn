package user

import (
	"fmt"
	"gin_learn/webook/service"
	regexp "github.com/dlclark/regexp2"
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

	fmt.Printf("%v", req)
	ctx.String(http.StatusOK, "注册成功")
	return
}
func (user *UserHandler) Login(ctx *gin.Context) {
	println("login")
}

func (user *UserHandler) Edit(ctx *gin.Context) {

}
func (user *UserHandler) Profile(ctx *gin.Context) {

}
