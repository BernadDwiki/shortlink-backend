package router

import (
	"github.com/BernadDwiki/shortlink-backend/internal/controller"
	"github.com/BernadDwiki/shortlink-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authController *controller.AuthController, linkController *controller.LinkController, jwtSecret string) {
	api := r.Group("/api")

	routes := api.Group("/auth")
	RegisterAuthRoutes(routes, authController)

	protected := api.Group("/links")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		protected.POST("", linkController.CreateLink)
		protected.GET("", linkController.GetUserLinks)
		protected.DELETE("/:id", linkController.DeleteLink)
	}

	// Public redirect endpoint
	r.GET("/:slug", linkController.GetLink)
}
