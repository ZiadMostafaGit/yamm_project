package handler

import (
	"net/http"
	"os"
	"strings"
	"yamm-project/app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is requeired"})
			c.Abort()
			return

		}

		partes := strings.Split(authHeader, " ")
		if len(partes) != 2 || partes[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is requeired"})
			c.Abort()
			return

		}

		jwtTocken := partes[1]
		secret := os.Getenv("JWT_SECRET")

		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(jwtTocken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("role", claims.Role)
		c.Set("store_id", claims.StoreID)

		c.Next()
	}

}
