package middleware

import (
	"net/http"
	"strings"

	"github.com/BernadDwiki/shortlink-backend/internal/jwt"
	"github.com/BernadDwiki/shortlink-backend/internal/response"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			response.Error(ctx, http.StatusUnauthorized, "missing authorization header")
			ctx.Abort()
			return
		}

		parts := strings.Fields(authorization)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(ctx, http.StatusUnauthorized, "invalid authorization header")
			ctx.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1], secret)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "invalid or expired token")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Next()
	}
}
