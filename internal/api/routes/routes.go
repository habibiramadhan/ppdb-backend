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
	// Inisialisasi repositories
	userRepo := repositories.NewUserRepository(cfg.DB)
	verificationRepo := repositories.NewVerificationRepository(cfg.DB)
	passwordResetRepo := repositories.NewPasswordResetRepository(cfg.DB)
	academicYearRepo := repositories.NewAcademicYearRepository(cfg.DB)
	majorRepo := repositories.NewMajorRepository(cfg.DB)
	majorFileRepo := repositories.NewMajorFileRepository(cfg.DB)
	quotaRepo := repositories.NewMajorQuotaRepository(cfg.DB)
	scheduleRepo := repositories.NewScheduleRepository(cfg.DB)
	notificationRepo := repositories.NewScheduleNotificationRepository(cfg.DB)

	// Setup email service
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
	academicYearService := services.NewAcademicYearService(academicYearRepo)
	majorService := services.NewMajorService(majorRepo, majorFileRepo)
	quotaService := services.NewMajorQuotaService(quotaRepo, academicYearRepo, majorRepo)
	scheduleService := services.NewScheduleService(scheduleRepo, notificationRepo, academicYearRepo, emailService)

	// Inisialisasi handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	verificationHandler := handlers.NewVerificationHandler(verificationService)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)
	academicYearHandler := handlers.NewAcademicYearHandler(academicYearService)
	majorHandler := handlers.NewMajorHandler(majorService)
	quotaHandler := handlers.NewMajorQuotaHandler(quotaService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)

	// Setup session
	sessionRepo := repositories.NewSessionRepository(cfg.DB)
	sessionService := services.NewSessionService(sessionRepo, userRepo)
	sessionHandler := handlers.NewSessionHandler(sessionService)

	// Route publik
	public := e.Group("/api/v1")

	// Route autentikasi
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)

	// Route verifikasi email
	public.GET("/verify-email", verificationHandler.VerifyEmail)
	public.POST("/resend-verification", verificationHandler.ResendVerification)

	// Route reset password
	public.POST("/forgot-password", passwordResetHandler.RequestReset)
	public.GET("/validate-reset-token", passwordResetHandler.ValidateToken)
	public.POST("/reset-password", passwordResetHandler.ResetPassword)

	// Route yang butuh autentikasi
	protected := e.Group("/api/v1")
	protected.Use(middlewares.JWTMiddleware(authService))

	protected.GET("/user/profile", authHandler.GetProfile)

	// Route khusus admin
	admin := protected.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	// Route manajemen user
	admin.GET("/users", adminHandler.GetAllUsers)
	admin.GET("/users/:id", adminHandler.GetUserByID)
	admin.PUT("/users/:id", adminHandler.UpdateUser)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)
	admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)

	// Route manajemen tahun akademik
	admin.POST("/academic-years", academicYearHandler.Create)
	admin.GET("/academic-years", academicYearHandler.GetAll) 
	admin.GET("/academic-years/:id", academicYearHandler.GetByID)
	admin.PUT("/academic-years/:id", academicYearHandler.Update)
	admin.DELETE("/academic-years/:id", academicYearHandler.Delete)
	admin.PATCH("/academic-years/:id/status", academicYearHandler.SetStatus)

	// Route manajemen jurusan
	admin.POST("/majors", majorHandler.Create)
	admin.PUT("/majors/:id", majorHandler.Update)
	admin.DELETE("/majors/:id", majorHandler.Delete)
	admin.PATCH("/majors/:id/status", majorHandler.SetStatus)
	admin.POST("/majors/:id/icon", majorHandler.UploadIcon)
	admin.POST("/majors/:id/files", majorHandler.UploadFiles)
	admin.DELETE("/majors/files/:id", majorHandler.DeleteFile)

	// Route manajemen kuota jurusan
	admin.POST("/major-quotas", quotaHandler.Create)
	admin.PUT("/major-quotas/:id", quotaHandler.Update)
	admin.DELETE("/major-quotas/:id", quotaHandler.Delete)
	admin.GET("/major-quotas/:id/logs", quotaHandler.GetLogs)

	// Route manajemen jadwal
	admin.POST("/schedules", scheduleHandler.Create)
	admin.PUT("/schedules/:id", scheduleHandler.Update)
	admin.DELETE("/schedules/:id", scheduleHandler.Delete)
	admin.PATCH("/schedules/:id/status", scheduleHandler.SetStatus)

	// Route manajemen sesi
	protected.GET("/sessions", sessionHandler.GetActiveSessions)
	protected.DELETE("/sessions/:id", sessionHandler.RevokeSession)
	protected.DELETE("/sessions", sessionHandler.RevokeAllSessions)

	// Route publik
	public.GET("/academic-years/active", academicYearHandler.GetActive)

	public.GET("/majors", majorHandler.GetAll)
	public.GET("/majors/:id", majorHandler.GetByID) 
	public.GET("/majors/:id/files", majorHandler.GetFiles)

	public.GET("/major-quotas", quotaHandler.GetAll)
	public.GET("/major-quotas/:id", quotaHandler.GetByID)

	public.GET("/schedules", scheduleHandler.GetAll)
	public.GET("/schedules/upcoming", scheduleHandler.GetUpcoming)
	public.GET("/schedules/academic-year/:yearId", scheduleHandler.GetByAcademicYear)
	public.GET("/schedules/:id", scheduleHandler.GetByID)
}
