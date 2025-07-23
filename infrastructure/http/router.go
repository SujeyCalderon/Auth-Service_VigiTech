package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"auth-service/application"
)

func ConfigureCORS(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func RegisterRoutes(router *gin.Engine, authService *application.Service) {
	api := router.Group("/api")
	{
		api.OPTIONS("/*any", func(c *gin.Context) { c.Status(204) })
		api.POST("/register", RegisterUserHandler(authService))
		api.POST("/login", LoginUserHandler(authService))
		api.GET("/users", GetAllUsersHandler(authService))
		api.PUT("/users/:id/status", UpdateUserStatusHandler(authService))
	}
}
