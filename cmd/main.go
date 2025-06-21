package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"moon/internal/config"
	"moon/internal/database"
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

		// TODO: Add more routes here
		// api.POST("/auth/register", authHandler.Register)
		// api.POST("/auth/login", authHandler.Login)
		//
		// // Protected routes
		// protected := api.Group("/")
		// protected.Use(middleware.AuthMiddleware())
		// {
		//     protected.GET("/users/profile", userHandler.GetProfile)
		//     protected.PUT("/users/profile", userHandler.UpdateProfile)
		// }
	}

	return r
}
