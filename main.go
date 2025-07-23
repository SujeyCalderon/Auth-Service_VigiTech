package main

import (
	"log"
	netHttp "net/http"

	"github.com/gin-gonic/gin"

	httpHandler "auth-service/infrastructure/http"
	"auth-service/infrastructure/db"
	"auth-service/infrastructure/token"
	"auth-service/helpers"
	"auth-service/application"
)

func main() {
	dbConn, err := helpers.ConnectMySQL()
	if err != nil {
		log.Fatalf("Error conectando a MySQL: %v", err)
	}
	defer dbConn.Close()

	userRepo := db.NewUserRepository(dbConn)
	jwtService := token.NewJWTService()
	authService := application.NewService(userRepo, jwtService)

	router := gin.Default()
	httpHandler.ConfigureCORS(router)

	//Rutas 
	public := router.Group("/api")
	{
		public.POST("/register", httpHandler.RegisterUserHandler(authService))
		public.POST("/login",    httpHandler.LoginUserHandler(authService))
	}

	protected := router.Group("/api")
	protected.Use(httpHandler.JWTAuthMiddleware(jwtService))
	{
		protected.GET("/users",            httpHandler.GetAllUsersHandler(authService))
		protected.PUT("/users/:id/status", httpHandler.UpdateUserStatusHandler(authService))
	}

	log.Println("Servidor iniciado en :8080")
	log.Fatal(netHttp.ListenAndServe(":8080", router))
}
