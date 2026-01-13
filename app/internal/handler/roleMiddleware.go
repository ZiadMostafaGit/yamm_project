package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		requestRole, exsist := c.Get("role")
		if !exsist {
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
			SendResponse(c, http.StatusForbidden, "You do not have permission to perform this action", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
