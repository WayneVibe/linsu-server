package router

import (
	"github.com/gin-gonic/gin"

	"app/internal/auths"
)

type AuthRouter struct {
}

// Register 负责注册用户相关的路由
func (u *AuthRouter) Register(engine *gin.Engine) {
	// 创建一个路由组
	userGroup := engine.Group("/api/v1/auth")
	{
		userHandler := auths.NewHandler()
		userGroup.POST("/register", userHandler.Register)
		userGroup.GET("/verify-email", userHandler.VerifyEmail)
		userGroup.POST("/login", userHandler.Login)
	}
}
