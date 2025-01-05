// internal/api/routes/routes.go
// file untuk setup routing aplikasi
package routes

import (
	"log"
	"os"
	"ppdb-backend/config"
	"ppdb-backend/internal/api/handlers"
	"ppdb-backend/internal/api/middlewares"
	"ppdb-backend/internal/core/repositories"
	"ppdb-backend/internal/core/services"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, cfg *config.Config) {
	userRepo := repositories.NewUserRepository(cfg.DB)
	verificationRepo := repositories.NewVerificationRepository(cfg.DB)

	emailService, err := services.NewEmailService()
	if err != nil {
		log.Fatal("Waduh gagal inisialisasi service email nih:", err)
	}

	authService := services.NewAuthService(
		userRepo,
		verificationRepo,
		emailService,
		os.Getenv("JWT_SECRET"),
	)
	adminService := services.NewAdminService(userRepo)
	verificationService := services.NewVerificationService(
		verificationRepo,
		userRepo,
		emailService,
	)
    
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	verificationHandler := handlers.NewVerificationHandler(verificationService)

	// Grup route yang bisa diakses publik
	public := e.Group("/api/v1")

	// Route buat autentikasi
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)

	// Route buat verifikasi
	public.GET("/verify-email", verificationHandler.VerifyEmail)
	public.POST("/resend-verification", verificationHandler.ResendVerification)

	// Grup route yang perlu login dulu
	protected := e.Group("/api/v1")
	protected.Use(middlewares.JWTMiddleware(authService))

	// Route yang butuh login
	protected.GET("/user/profile", authHandler.GetProfile)

	// Grup route khusus admin
	admin := protected.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	// Route buat admin
	admin.GET("/users", adminHandler.GetAllUsers)
	admin.GET("/users/:id", adminHandler.GetUserByID)
	admin.PUT("/users/:id", adminHandler.UpdateUser)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)
	admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)
}