// internal/api/routes/routes.go
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
	passwordResetRepo := repositories.NewPasswordResetRepository(cfg.DB)
	academicYearRepo := repositories.NewAcademicYearRepository(cfg.DB)

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
	passwordResetService := services.NewPasswordResetService(
		passwordResetRepo,
		userRepo,
		emailService,
	)
	academicYearService := services.NewAcademicYearService(academicYearRepo)

	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	verificationHandler := handlers.NewVerificationHandler(verificationService)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)
	academicYearHandler := handlers.NewAcademicYearHandler(academicYearService)

	sessionRepo := repositories.NewSessionRepository(cfg.DB)
	sessionService := services.NewSessionService(sessionRepo, userRepo)
	sessionHandler := handlers.NewSessionHandler(sessionService)

	public := e.Group("/api/v1")

	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)

	public.GET("/verify-email", verificationHandler.VerifyEmail)
	public.POST("/resend-verification", verificationHandler.ResendVerification)

	public.POST("/forgot-password", passwordResetHandler.RequestReset)
	public.GET("/validate-reset-token", passwordResetHandler.ValidateToken)
	public.POST("/reset-password", passwordResetHandler.ResetPassword)

	protected := e.Group("/api/v1")
	protected.Use(middlewares.JWTMiddleware(authService))

	protected.GET("/user/profile", authHandler.GetProfile)

	admin := protected.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	admin.GET("/users", adminHandler.GetAllUsers)
	admin.GET("/users/:id", adminHandler.GetUserByID)
	admin.PUT("/users/:id", adminHandler.UpdateUser)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)
	admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)

	admin.POST("/academic-years", academicYearHandler.Create)
	admin.GET("/academic-years", academicYearHandler.GetAll) 
	admin.GET("/academic-years/:id", academicYearHandler.GetByID)
	admin.PUT("/academic-years/:id", academicYearHandler.Update)
	admin.DELETE("/academic-years/:id", academicYearHandler.Delete)
	admin.PATCH("/academic-years/:id/status", academicYearHandler.SetStatus)

	

	protected.GET("/sessions", sessionHandler.GetActiveSessions)
	protected.DELETE("/sessions/:id", sessionHandler.RevokeSession)
	protected.DELETE("/sessions", sessionHandler.RevokeAllSessions)

	// Route publik
	public.GET("/academic-years/active", academicYearHandler.GetActive)

}
