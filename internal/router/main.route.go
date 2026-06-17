package router

import (
	"github.com/BernadDwiki/shortlink-backend/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authController *controller.AuthController, linkController *controller.LinkController, jwtSecret string) {
	api := r.Group("/api")

	// Auth routes
	RegisterAuthRoutes(api, authController)

	// Links routes (protected + public)
	RegisterLinkRoutes(api, linkController, jwtSecret)
	RegisterPublicRoutes(r, linkController)
}
