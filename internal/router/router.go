package router

import (
	"golang-api-template/internal/config"
	"golang-api-template/internal/handlers"
	"golang-api-template/internal/i18n"
	"golang-api-template/internal/middlewares"
	"golang-api-template/internal/repository"
	"golang-api-template/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *gin.Engine {
	if err := i18n.Initialize(); err != nil {
		panic("Failed to load translations: " + err.Error())
	}

	r := gin.Default()

	// Middlewares
	r.Use(middlewares.LocaleMiddleware())

	// Repos
	userRepo := repository.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepo) // from previous examples
	authService := service.NewAuthService(userRepo, rdb, cfg)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	v1 := r.Group("/api/v1")
	{
		// Auth
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/refresh", authHandler.RefreshToken)
		v1.POST("/auth/logout", authHandler.Logout)

		// Registration - example
		// v1.POST("/users/register", authHandler.Register)
	}

	// Protected routes
	auth := v1.Group("/users")
	auth.Use(middlewares.AuthMiddleware(cfg)) // e.g. checks valid JWT
	{
		auth.GET("/:id", userHandler.GetByID)
		auth.GET("/", userHandler.List)
		auth.PUT("/:id", userHandler.Update)
		auth.DELETE("/:id", userHandler.Delete)
	}

	return r
}
