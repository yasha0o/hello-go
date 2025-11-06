package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Code        int `json:"code"`
	Description any `json:"description"`
}

type ValidationError struct {
	Err error
}

func (v *ValidationError) Error() string {
	return v.Err.Error()
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
	var error ApiError

	switch e := err.(type) {
	case *ValidationError:
		error = ApiError{
			Code:        http.StatusBadRequest,
			Description: e.Error(),
		}

	default:
		error = ApiError{
			Code:        http.StatusInternalServerError,
			Description: e.Error(),
		}
	}

	c.JSON(error.Code, error)
}
