package router

import (
	"github.com/BernadDwiki/shortlink-backend/internal/controller"
	"github.com/BernadDwiki/shortlink-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterLinkRoutes registers protected link-related API routes under the provided router group.
func RegisterLinkRoutes(rg *gin.RouterGroup, linkController *controller.LinkController, jwtSecret string) {
	protected := rg.Group("/links")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		protected.POST("", linkController.CreateLink)
		protected.GET("", linkController.GetUserLinks)
		protected.DELETE("/:id", linkController.DeleteLink)
	}
}

// RegisterPublicRoutes registers public, unauthenticated routes (like redirect) on the engine.
func RegisterPublicRoutes(r *gin.Engine, linkController *controller.LinkController) {
	r.GET("/:slug", linkController.GetLink)
}
