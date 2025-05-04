// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	IsError bool        `json:"isError"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "success"
	}

	c.JSON(http.StatusOK, Response{
		IsError: false,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	if message == "" {
		message = "error"
	}

	c.JSON(statusCode, Response{
		IsError: true,
		Message: message,
		Data:    nil,
	})
}
