package main

import (
	"yamm-project/app/internal/config"
	"yamm-project/app/internal/db"
	"yamm-project/app/internal/handler"
	"yamm-project/app/internal/repository"
	"yamm-project/app/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		println(err)
	}

	userRepo := repository.NewUserRepository(database)
	storeRepo := repository.NewStoreRepo(database)
	authService := service.NewAuthService(userRepo, storeRepo, database)

	authHandler := handler.NewAuthHandler(authService)
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/registerNewAdminForYammApp", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	r.Run(":8080")
	println("server is running on port 8080")

}
