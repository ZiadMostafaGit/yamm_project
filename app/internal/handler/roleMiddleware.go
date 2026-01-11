package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		requestRole, exsist := c.Get("role")
		if !exsist {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found in context"})
			c.Abort()
			return
		}
		isAllowedToContinue := false
		for _, role := range allowedRoles {
			if role == requestRole {
				isAllowedToContinue = true
				break
			}
		}

		if !isAllowedToContinue {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
			c.Abort()
			return
		}

		c.Next()

	}

}
