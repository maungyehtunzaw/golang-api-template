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

	emailConfig := config.GetEmailConfig()
	emailService := service.NewEmailService(emailConfig)

	// Repos
	userRepo := repository.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepo) // from previous examples
	authService := service.NewAuthService(userRepo, rdb, cfg)

	// Handlers
	userHandler := handlers.NewUserHandler(userService, emailService)
	authHandler := handlers.NewAuthHandler(authService)

	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handlers.NewRoleHandler(roleService)

	// Public routes
	v1 := r.Group("/api/v1")
	{
		// Auth
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/refresh", authHandler.RefreshToken)
		v1.POST("/auth/logout", authHandler.Logout)
		v1.POST("/auth/register", userHandler.Create)
		v1.POST("/auth/forgot-password", userHandler.ForgotPassword)

		v1.POST("/roles", roleHandler.CreateRole)
		v1.GET("/roles", roleHandler.GetAllRoles)
		v1.GET("/roles/:id", roleHandler.GetRoleByID)
		v1.PUT("/roles/:id", roleHandler.UpdateRole)
		v1.DELETE("/roles/:id", roleHandler.DeleteRole)

		v1.GET("/roles/:id/permissions", roleHandler.GetPermissionsByRoleID)
		v1.GET("/users/:id/permissions", userHandler.GetPermissionsByUserID)
	}

	// Protected routes
	auth := v1.Group("/users")

	auth.Use(middlewares.AuthMiddleware(cfg)) // e.g. checks valid JWT
	{
		auth.GET("/:id", userHandler.GetByID)
		auth.GET("/", userHandler.List)
		auth.GET("/getuser", authHandler.GetAuthUser)

		auth.PUT("/:id", userHandler.Update)
		auth.DELETE("/:id", userHandler.Delete)
	}

	return r
}
