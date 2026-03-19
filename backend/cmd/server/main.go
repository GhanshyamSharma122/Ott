package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shyam/ott-backend/internal/config"
	"github.com/shyam/ott-backend/internal/database"
	"github.com/shyam/ott-backend/internal/handlers"
	"github.com/shyam/ott-backend/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db := database.Connect(cfg)

	// Run migrations
	database.RunMigrations(db, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg.JWTSecret)
	videoHandler := handlers.NewVideoHandler()
	categoryHandler := handlers.NewCategoryHandler()
	watchlistHandler := handlers.NewWatchlistHandler()
	historyHandler := handlers.NewHistoryHandler()
	userHandler := handlers.NewUserHandler()
	paymentHandler := handlers.NewPaymentHandler()

	// Setup Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	// Authenticated routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Auth
		protected.GET("/auth/me", authHandler.Me)

		// Videos
		protected.GET("/videos", videoHandler.ListVideos)
		protected.GET("/videos/continue", videoHandler.ContinueWatching)
		protected.GET("/videos/:id", videoHandler.GetVideo)

		// Categories
		protected.GET("/categories", categoryHandler.ListCategories)
		protected.GET("/categories/:id", categoryHandler.GetCategory)

		// Watchlist
		protected.GET("/watchlist", watchlistHandler.GetWatchlist)
		protected.POST("/watchlist", watchlistHandler.AddToWatchlist)
		protected.DELETE("/watchlist/:video_id", watchlistHandler.RemoveFromWatchlist)
		protected.GET("/watchlist/check/:video_id", watchlistHandler.CheckWatchlist)

		// Watch History
		protected.GET("/history", historyHandler.GetHistory)
		protected.POST("/history", historyHandler.UpdateProgress)

		// Profile
		protected.GET("/profile", userHandler.GetProfile)
		protected.PUT("/profile", userHandler.UpdateProfile)

		// Payment
		protected.POST("/payment/subscribe", paymentHandler.Subscribe)
		protected.GET("/payment/status", paymentHandler.GetSubscriptionStatus)
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(middleware.RequireAdmin())
	{
		// Videos
		admin.GET("/videos", videoHandler.AdminListVideos)
		admin.POST("/videos", videoHandler.AdminCreateVideo)
		admin.PUT("/videos/:id", videoHandler.AdminUpdateVideo)
		admin.DELETE("/videos/:id", videoHandler.AdminDeleteVideo)

		// Categories
		admin.POST("/categories", categoryHandler.AdminCreateCategory)
		admin.PUT("/categories/:id", categoryHandler.AdminUpdateCategory)
		admin.DELETE("/categories/:id", categoryHandler.AdminDeleteCategory)

		// Users
		admin.GET("/users", userHandler.AdminListUsers)
		admin.PUT("/users/:id", userHandler.AdminUpdateUser)

		// Stats
		admin.GET("/stats", userHandler.AdminStats)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
