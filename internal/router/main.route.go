package router

import (
	"github.com/BernadDwiki/shortlink-backend/internal/controller"
	"github.com/BernadDwiki/shortlink-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authController *controller.AuthController, jwtSecret string) {
	api := r.Group("/api")

	routes := api.Group("/auth")
	RegisterAuthRoutes(routes, authController)

	protected := api.Group("/links")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		protected.GET("", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "protected links endpoint"})
		})
	}
}
