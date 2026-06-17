package router

import (
	"github.com/BernadDwiki/shortlink-backend/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authController *controller.AuthController) {
	rg.POST("/register", authController.Register)
	rg.POST("/login", authController.Login)
}
