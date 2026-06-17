package controller

import (
	"errors"
	"net/http"

	"github.com/BernadDwiki/shortlink-backend/internal/dto"
	"github.com/BernadDwiki/shortlink-backend/internal/response"
	"github.com/BernadDwiki/shortlink-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (a *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.authService.Register(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyRegistered) {
			response.Error(ctx, http.StatusBadRequest, "email already registered")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to register user")
		return
	}

	response.Success(ctx, http.StatusCreated, "registration successful", gin.H{"user": gin.H{"id": user.ID, "email": user.Email}})
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := a.authService.Login(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(ctx, http.StatusUnauthorized, "invalid credentials")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to create token")
		return
	}

	response.Success(ctx, http.StatusOK, "Login successful", gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email}})
}
