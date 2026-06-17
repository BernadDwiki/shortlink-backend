package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Results interface{} `json:"results,omitempty"`
}

func Success(ctx *gin.Context, status int, message string, results interface{}) {
	ctx.JSON(status, APIResponse{Success: true, Message: message, Results: results})
}

func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, APIResponse{Success: false, Message: message})
}
