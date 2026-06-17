package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/BernadDwiki/shortlink-backend/internal/config"
	"github.com/BernadDwiki/shortlink-backend/internal/controller"
	"github.com/BernadDwiki/shortlink-backend/internal/middleware"
	"github.com/BernadDwiki/shortlink-backend/internal/repository"
	"github.com/BernadDwiki/shortlink-backend/internal/router"
	"github.com/BernadDwiki/shortlink-backend/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

	repo := repository.NewUserRepository(db)
	authService := service.NewAuthService(repo, cfg.JWTSecret, time.Minute*time.Duration(cfg.JWTExpirationMinutes))
	authController := controller.NewAuthController(authService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware)

	router.RegisterRoutes(r, authController, cfg.JWTSecret)

	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	log.Printf("starting server at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
