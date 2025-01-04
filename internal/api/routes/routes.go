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

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)

    // Public routes group
    public := e.Group("/api/v1")
    
    // Auth routes
    public.POST("/auth/register", authHandler.Register)
    public.POST("/auth/login", authHandler.Login)

    // Protected routes group
    protected := e.Group("/api/v1")
    protected.Use(middlewares.JWTMiddleware(authService))

    // Protected routes will be added here
    protected.GET("/user/profile", authHandler.GetProfile)
}