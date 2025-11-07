package controller

import (
	"log"
	"net/http"
	"runtime/debug"

	"hello-go/internal/domain"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Code        int `json:"code"`
	Description any `json:"description"`
}

type HandlerFunc func(c *gin.Context) error

func ErrorWrapper(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	log.Printf("Request error: %v\nStack trace:\n%s", err, debug.Stack())

	var apiErr ApiError

	switch e := err.(type) {
	case *domain.ValidationError:
		apiErr = ApiError{
			Code:        http.StatusBadRequest,
			Description: e.Error(),
		}
	case *domain.NotFoundError:
		apiErr = ApiError{
			Code:        http.StatusNotFound,
			Description: e.Error(),
		}
	default:
		apiErr = ApiError{
			Code:        http.StatusInternalServerError,
			Description: e.Error(),
		}
	}

	c.JSON(apiErr.Code, apiErr)
}
