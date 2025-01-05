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
    // Inisialisasi repositories
    userRepo := repositories.NewUserRepository(cfg.DB)
    verificationRepo := repositories.NewVerificationRepository(cfg.DB)
    passwordResetRepo := repositories.NewPasswordResetRepository(cfg.DB)

    // Inisialisasi email service
    emailService, err := services.NewEmailService()
    if err != nil {
        log.Fatal("Waduh gagal inisialisasi service email nih:", err)
    }

    // Inisialisasi services
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

    // Inisialisasi handlers
    authHandler := handlers.NewAuthHandler(authService)
    adminHandler := handlers.NewAdminHandler(adminService)
    verificationHandler := handlers.NewVerificationHandler(verificationService)
    passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)

    // Grup route yang bisa diakses publik
    public := e.Group("/api/v1")

    // Route untuk autentikasi
    public.POST("/auth/register", authHandler.Register)            // Daftar user baru
    public.POST("/auth/login", authHandler.Login)                  // Login user

    // Route untuk verifikasi email
    public.GET("/verify-email", verificationHandler.VerifyEmail)           // Verifikasi email dari link
    public.POST("/resend-verification", verificationHandler.ResendVerification) // Kirim ulang email verifikasi

    // Route untuk reset password
    public.POST("/forgot-password", passwordResetHandler.RequestReset)     // Request reset password
    public.GET("/validate-reset-token", passwordResetHandler.ValidateToken)// Validasi token reset
    public.POST("/reset-password", passwordResetHandler.ResetPassword)     // Reset password dengan token

    // Grup route yang perlu login (protected)
    protected := e.Group("/api/v1")
    protected.Use(middlewares.JWTMiddleware(authService))

    // Route yang butuh login
    protected.GET("/user/profile", authHandler.GetProfile)                 // Get profil user

    // Grup route khusus admin
    admin := protected.Group("/admin")
    admin.Use(middlewares.AdminMiddleware())

    // Route untuk manajemen user oleh admin
    admin.GET("/users", adminHandler.GetAllUsers)                         // List semua user
    admin.GET("/users/:id", adminHandler.GetUserByID)                    // Get detail user
    admin.PUT("/users/:id", adminHandler.UpdateUser)                     // Update data user
    admin.DELETE("/users/:id", adminHandler.DeleteUser)                  // Hapus user
    admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)      // Update status user
}