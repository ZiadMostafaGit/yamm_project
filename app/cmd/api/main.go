package main

import (
	"net/http"
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
	categoryRepo := repository.NewCategoryRepository(database)

	faqRepo := repository.NewFAQRepository(database)
	transRepo := repository.NewTranslationRepository(database)

	authService := service.NewAuthService(userRepo, storeRepo, database)
	categoryService := service.NewCategoryService(categoryRepo)

	faqService := service.NewFAQService(faqRepo, transRepo, categoryRepo, storeRepo, database)

	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	faqHandler := handler.NewFAQHandler(faqService)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/registerNewAdminForYammApp", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		app := api.Group("/")
		{
			app.GET("/customerPage", handler.AuthMiddleware(), handler.RoleMiddleware("customer", "merchant"), func(ctx *gin.Context) {
				ctx.JSON(http.StatusAccepted, gin.H{"middleware auth": "jwt middleware is ok and role is accepted"})
			})
			app.GET("/adminPage", handler.AuthMiddleware(), handler.RoleMiddleware("admin"), func(ctx *gin.Context) {
				ctx.JSON(http.StatusAccepted, gin.H{"middleware auth": "jwt middleware is ok and role is accepted"})
			})
		}

		category := api.Group("/category")
		{
			category.POST("/Create", handler.AuthMiddleware(), handler.RoleMiddleware("admin"), categoryHandler.Create)
			category.PATCH("/Update/:id", handler.AuthMiddleware(), handler.RoleMiddleware("admin"), categoryHandler.Update)
			category.DELETE("/Delete/:id", handler.AuthMiddleware(), handler.RoleMiddleware("admin"), categoryHandler.Delete)
			category.GET("/GetAll", handler.AuthMiddleware(), handler.RoleMiddleware("customer", "merchant", "admin"), categoryHandler.GetAll)
		}

		faq := api.Group("/faq")
		{
			faq.GET("/GetForCustomer", faqHandler.GetCustomerView)

			faq.POST("/Create", handler.AuthMiddleware(), handler.RoleMiddleware("admin", "merchant"), faqHandler.Create)

			faq.PUT("/UpdateTranslations/:id", handler.AuthMiddleware(), handler.RoleMiddleware("admin", "merchant"), faqHandler.UpdateTranslations)

			faq.DELETE("/Delete/:id", handler.AuthMiddleware(), handler.RoleMiddleware("admin", "merchant"), faqHandler.Delete)

			faq.PATCH("/UpdateVisibility/:id", handler.AuthMiddleware(), handler.RoleMiddleware("admin"), faqHandler.UpdateVisibility)
		}
	}

	r.Run(":8080")
	println("server is running on port 8080")
}
