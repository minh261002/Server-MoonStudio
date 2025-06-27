package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"moon/internal/config"
	"moon/internal/database"
	"moon/internal/domain/user"
	httpHandler "moon/internal/handler/http"
	"moon/internal/middleware"
	"moon/internal/repository"
	"moon/internal/usecase"
	"moon/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	if err := config.LoadConfig("configs/config.yaml"); err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	cfg := config.GetConfig()

	// Initialize logger
	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	log := logger.GetLogger()
	log.Info("Starting Moon API", zap.String("version", cfg.App.Version))

	// Set Gin mode
	gin.SetMode(cfg.App.Mode)

	// Connect to database
	if err := database.ConnectDatabase(cfg); err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	log.Info("Connected to database successfully")

	// Auto migrate
	db := database.GetDB()
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}
	log.Info("Database migration completed")

	// Setup router
	r := setupRouter()

	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		log.Info("Server starting", zap.String("address", addr))
		if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Close database connection
	if err := database.CloseDatabase(); err != nil {
		log.Error("Error closing database", zap.Error(err))
	}

	log.Info("Server exited")
}

func setupRouter() *gin.Engine {
	cfg := config.GetConfig()
	db := database.GetDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg)

	// Initialize handlers
	authHandler := httpHandler.NewAuthHandler(authUseCase)

	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"version": "1.0.0",
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile routes
			protected.GET("/profile", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				email, _ := c.Get("email")
				role, _ := c.Get("role")

				c.JSON(http.StatusOK, gin.H{
					"user_id": userID,
					"email":   email,
					"role":    role,
				})
			})
			// TODO: Add more protected routes here
		}
	}

	return r
}
