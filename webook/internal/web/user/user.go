package user

import (
	"errors"
	"fmt"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository/cache"
	"gin_learn/webook/internal/service"
	"gin_learn/webook/internal/web/utils"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// UserHandler 定义所有有关user的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	nicknameExp *regexp.Regexp
	birthdayExp *regexp.Regexp
	profileExp  *regexp.Regexp
	passwordExp *regexp.Regexp
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 生命要放进token里面的数据
	Uid       int64
	UserAgent string
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	nicknameExp := regexp.MustCompile(nicknameRegexPattern, regexp.None)
	birthdayExp := regexp.MustCompile(birthdayRegexPattern, regexp.None)
	profileExp := regexp.MustCompile(profileRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		nicknameExp: nicknameExp,
		birthdayExp: birthdayExp,
		profileExp:  profileExp,
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
	if errors.Is(err, service.ErrUserDuplicate) {
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

func (user *UserHandler) LoginJWT(ctx *gin.Context) {
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
	//sess := sessions.Default(ctx)
	//// 需要放在session里面的值
	//sess.Set("userId", findUser.Id)
	//err = sess.Save()
	//if err != nil {
	//	return
	//}
	if err = user.setJWTToken(ctx, findUser.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{"msg": "系统错误", "success": false})
		return
	}
	fmt.Println(findUser)
	ctx.JSON(http.StatusOK, map[string]any{"msg": "登录成功", "success": true})
	return
}

func (user *UserHandler) setJWTToken(ctx *gin.Context, uId int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		UserAgent: ctx.Request.UserAgent(),
		Uid:       uId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("r4BKnmqBgWhnudRc4xufW9f97ODTqX10"))
	if err != nil {

		return err
	}
	fmt.Println(tokenStr)
	ctx.Header("x-jwt-token", tokenStr)
	return nil
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
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		Profile  string `json:"profile"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验nickname
	ok, err := user.nicknameExp.MatchString(req.Nickname)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "修改失败, 昵称长度应在1-64个字节", "success": false})
		return
	}
	// 校验birthday
	ok, err = user.birthdayExp.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "修改失败, 生日日期格式错误", "success": false})
		return
	}

	// 校验profile
	ok, err = user.profileExp.MatchString(req.Profile)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "修改失败,个人简介长度应在256个字节内", "success": false})
		return
	}

	// 获取userId
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		// 监控输出
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	UserId := claims.Uid
	//userId := claims.Uid
	//session := sessions.Default(ctx)
	//UserId := session.Get("userId").(int64)
	err = user.svc.Edit(ctx, domain.UserProfile{
		UserId:   UserId,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		Profile:  req.Profile,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "修改失败", "success": false})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"msg": "修改成功", "success": true})
	return
}
func (user *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	userId := sess.Get("userId").(int64)
	userProfile, err := user.svc.Profile(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "获取个人信息失败", "success": false})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"msg": "", "success": true, "data": map[string]any{
		"nickname": userProfile.Nickname,
		"profile":  userProfile.Profile,
		"birthday": userProfile.Birthday,
		"userId":   userProfile.UserId,
	}})
	return
}
func (user *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		// 监控输出
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	userId := claims.Uid
	userProfile, err := user.svc.Profile(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "获取个人信息失败", "success": false})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"msg": "", "success": true, "data": map[string]any{
		"nickname": userProfile.Nickname,
		"profile":  userProfile.Profile,
		"birthday": userProfile.Birthday,
		"userId":   userProfile.UserId,
	}})
	return
}

func (user *UserHandler) SendCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	const biz = "login"
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := user.codeSvc.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]any{"msg": "发送失败", "success": false})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"msg": "发送成功", "success": true})
	return
}

func (user *UserHandler) LoginSms(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"Code"`
	}
	const biz = "login"
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	res, err := user.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if errors.Is(err, cache.ErrCodeVerifyTooMany) {
		ctx.JSON(http.StatusOK, utils.Result{Code: 500, Success: false, Msg: "验证码校验过于频繁", Data: []string{}})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 500, Success: false, Msg: "系统错误", Data: []string{}})
		return
	}
	if !res {
		ctx.JSON(http.StatusOK, utils.Result{Code: 400, Success: false, Msg: "验证码错误，请重试", Data: []string{}})
		return
	}
	// 验证成功后的逻辑
	findUser, _err := user.svc.FindOrCreate(ctx, req.Phone)
	if _err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 500, Success: false, Msg: "系统错误", Data: []string{}})
		return
	}
	_err = user.setJWTToken(ctx, findUser.Id)
	if _err != nil {
		ctx.JSON(http.StatusOK, utils.Result{Code: 500, Success: false, Msg: "系统错误", Data: []string{}})
		return
	}
	ctx.JSON(http.StatusOK, utils.Result{Code: 200, Success: true, Msg: "验证成功", Data: []string{}})
	return
}
