// @title Shortlink Backend API
// @version 1.0
// @description Shortlink service API with JWT authentication, link creation, listing, deletion, and redirect.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/BernadDwiki/shortlink-backend/docs"
	"github.com/BernadDwiki/shortlink-backend/internal/config"
	"github.com/BernadDwiki/shortlink-backend/internal/controller"
	"github.com/BernadDwiki/shortlink-backend/internal/middleware"
	"github.com/BernadDwiki/shortlink-backend/internal/repository"
	"github.com/BernadDwiki/shortlink-backend/internal/router"
	"github.com/BernadDwiki/shortlink-backend/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, time.Minute*time.Duration(cfg.JWTExpirationMinutes))
	authController := controller.NewAuthController(authService)

	linkRepo := repository.NewLinkRepository(db)
	linkService := service.NewLinkService(linkRepo, cfg.BaseURL)
	linkController := controller.NewLinkController(linkService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.RegisterRoutes(r, authController, linkController, cfg.JWTSecret)

	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	log.Printf("starting server at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
