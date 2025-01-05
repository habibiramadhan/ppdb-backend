// internal/api/routes/routes.go
package routes

import (
	"os"
	"ppdb-backend/config"
	"ppdb-backend/internal/api/handlers"
	"ppdb-backend/internal/api/middlewares"
	"ppdb-backend/internal/core/repositories"
	"ppdb-backend/internal/core/services"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(cfg.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	adminService := services.NewAdminService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)

	// Public routes group
	public := e.Group("/api/v1")

	// Auth routes
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)

	// Protected routes group
	protected := e.Group("/api/v1")
	protected.Use(middlewares.JWTMiddleware(authService))

	// Protected routes
	protected.GET("/user/profile", authHandler.GetProfile)

	// Admin routes group
	admin := protected.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	// Admin routes
	admin.GET("/users", adminHandler.GetAllUsers)
	admin.GET("/users/:id", adminHandler.GetUserByID)
	admin.PUT("/users/:id", adminHandler.UpdateUser)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)
	admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)
}