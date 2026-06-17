package helper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validationFieldMap = map[string]string{
	"Email":       "email",
	"Password":    "password",
	"OriginalURL": "original_url",
	"Slug":        "slug",
}

func BindJSON(ctx *gin.Context, obj any) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return errors.New(formatValidationErrors(validationErrors))
		}
		return errors.New("invalid request payload")
	}

	return nil
}

func formatValidationErrors(errors validator.ValidationErrors) string {
	messages := make([]string, 0, len(errors))

	for _, fieldError := range errors {
		fieldName := fieldError.Field()
		if label, ok := validationFieldMap[fieldName]; ok {
			fieldName = label
		} else {
			fieldName = strings.ToLower(fieldName)
		}

		switch fieldError.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", fieldName))
		case "email":
			messages = append(messages, fmt.Sprintf("%s must be a valid email address", fieldName))
		case "url":
			messages = append(messages, fmt.Sprintf("%s must be a valid URL", fieldName))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s characters", fieldName, fieldError.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must be at most %s characters", fieldName, fieldError.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", fieldName))
		}
	}

	return strings.Join(messages, "; ")
}
