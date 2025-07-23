package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"auth-service/infrastructure/token"
)

func JWTAuthMiddleware(tokenService token.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Se requiere Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Formato de token inválido"})
			c.Abort()
			return
		}
		jwtToken := parts[1]

		claims, err := tokenService.ValidateToken(jwtToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token inválido: " + err.Error()})
			c.Abort()
			return
		}

		if sub, ok := claims["sub"].(string); ok {
			c.Set("userID", sub)
		}
		if r, ok := claims["role"].(string); ok {
			c.Set("role", r)
		}
		c.Next()
	}
}
