package user

import "github.com/gin-gonic/gin"

func (user *UserHandler) RegisterUserRouter(rg *gin.RouterGroup) {
	rg.POST("/signup", user.SignUp)
	rg.POST("/login", user.Login)
	rg.POST("/edit", user.Edit)
	rg.GET("/profile", user.Profile)
}
