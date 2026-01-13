package handler

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendResponse(c *gin.Context, statusCode int, message string, data any) {
	if data == nil {
		data = []string{}
	}

	c.JSON(statusCode, APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
	})
}
