package user

import "github.com/gin-gonic/gin"

func (user *UserHandler) RegisterUserRouter(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/signup", user.SignUp)
	ug.POST("/login", user.Login)
	ug.POST("/login_jwt", user.LoginJWT)
	ug.POST("/edit", user.Edit)
	//rg.GET("/profile", user.Profile)
	ug.GET("/profile", user.ProfileJWT)

	ug.POST("login_sms/code/send", user.SendCode)
	ug.POST("login_sms", user.LoginSms)
}
