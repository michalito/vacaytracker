package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/config"
	"vacaytracker-api/internal/handler"
	"vacaytracker-api/internal/middleware"
	"vacaytracker-api/internal/repository/sqlite"
	"vacaytracker-api/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Run database migrations
	if err := db.RunMigrations("./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := sqlite.NewUserRepository(db)
	vacationRepo := sqlite.NewVacationRepository(db)
	settingsRepo := sqlite.NewSettingsRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo)
	userService := service.NewUserService(userRepo, authService)
	emailService := service.NewEmailService(cfg)
	newsletterService := service.NewNewsletterService(cfg, userRepo, vacationRepo, settingsRepo, emailService)

	// Initialize and start the newsletter scheduler
	scheduler := service.NewScheduler(newsletterService, settingsRepo)
	scheduler.Start()

	// Create initial admin user if it doesn't exist
	if err := authService.CreateInitialAdmin(
		cfg.AdminEmail,
		cfg.AdminPassword,
		cfg.AdminName,
		25, // Default vacation balance
	); err != nil {
		log.Fatalf("Failed to create initial admin: %v", err)
	}

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	vacationHandler := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	adminHandler := handler.NewAdminHandler(cfg, userService, userRepo, vacationService, vacationRepo, settingsRepo, emailService, newsletterService)

	// Create Gin router
	router := gin.New()

	// Add global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorMiddleware())

	// Security headers middleware
	if cfg.IsProduction() {
		router.Use(middleware.ProductionSecurityHeaders())
	} else {
		router.Use(middleware.SecurityHeaders())
	}

	// Initialize security logger
	securityLogger := middleware.NewSecurityLogger()
	router.Use(middleware.SecurityLoggingMiddleware(securityLogger))

	// Initialize rate limiters
	loginRateLimiter := middleware.LoginRateLimiter()
	apiRateLimiter := middleware.APIRateLimiter()

	// CORS middleware (development mode allows all origins)
	if cfg.IsDevelopment() {
		router.Use(middleware.DefaultCORSMiddleware())
	} else {
		// In production, restrict to specific origins
		router.Use(middleware.CORSMiddleware([]string{cfg.AppURL}))
	}

	// Public routes
	router.GET("/health", healthHandler.Check)

	// API routes
	api := router.Group("/api")
	api.Use(apiRateLimiter.Middleware()) // Apply general rate limiting to all API routes
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			// Login has stricter rate limiting (5 per minute)
			auth.POST("/login", loginRateLimiter.Middleware(), authHandler.Login)
		}

		// Auth routes (authenticated)
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(authService))
		{
			authProtected.GET("/me", authHandler.Me)
			authProtected.PUT("/password", authHandler.ChangePassword)
			authProtected.PUT("/email-preferences", authHandler.UpdateEmailPreferences)
		}

		// Vacation routes (authenticated)
		vacation := api.Group("/vacation")
		vacation.Use(middleware.AuthMiddleware(authService))
		{
			vacation.POST("/request", vacationHandler.Create)
			vacation.GET("/requests", vacationHandler.List)
			vacation.GET("/requests/:id", vacationHandler.Get)
			vacation.DELETE("/requests/:id", vacationHandler.Cancel)
			vacation.GET("/team", vacationHandler.Team)
		}

		// Admin routes (authenticated + admin role)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		admin.Use(middleware.AdminMiddleware())
		{
			// User management
			admin.GET("/users", adminHandler.ListUsers)
			admin.POST("/users", adminHandler.CreateUser)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)
			admin.PUT("/users/:id/balance", adminHandler.UpdateBalance)
			admin.POST("/users/reset-balances", adminHandler.ResetBalances)

			// Vacation management
			admin.GET("/vacation/pending", adminHandler.ListPending)
			admin.PUT("/vacation/:id/review", adminHandler.Review)

			// Settings
			admin.GET("/settings", adminHandler.GetSettings)
			admin.PUT("/settings", adminHandler.UpdateSettings)

			// Newsletter
			admin.POST("/newsletter/send", adminHandler.SendNewsletter)
			admin.GET("/newsletter/preview", adminHandler.PreviewNewsletter)

			// Email Testing
			admin.POST("/email/test", adminHandler.SendTestEmail)
			admin.POST("/email/preview", adminHandler.PreviewEmail)
		}
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("VacayTracker API starting on port %s", cfg.Port)
		log.Printf("Environment: %s", cfg.Env)
		log.Printf("Health check: http://localhost:%s/health", cfg.Port)
		log.Printf("Login endpoint: POST http://localhost:%s/api/auth/login", cfg.Port)
		log.Printf("Vacation endpoints: http://localhost:%s/api/vacation/*", cfg.Port)
		log.Printf("Admin endpoints: http://localhost:%s/api/admin/*", cfg.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop the newsletter scheduler
	scheduler.Stop()

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
