package database

import (
	"log"

	"github.com/shyam/ott-backend/internal/config"
	"github.com/shyam/ott-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, cfg *config.Config) {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Video{},
		&models.WatchlistItem{},
		&models.WatchHistory{},
		&models.Payment{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")

	// Seed default data
	seedAdmin(db, cfg)
	seedCategories(db)
}

func seedCategories(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		log.Println("Categories already exist, skipping seed")
		return
	}

	categories := []models.Category{
		{Name: "Movies", Description: "Full-length feature films", SortOrder: 1},
		{Name: "TV Shows", Description: "Television series and episodes", SortOrder: 2},
		{Name: "Documentaries", Description: "Documentary films and series", SortOrder: 3},
		{Name: "Comedy", Description: "Comedy specials and shows", SortOrder: 4},
		{Name: "Action", Description: "Action and adventure content", SortOrder: 5},
		{Name: "Drama", Description: "Drama films and series", SortOrder: 6},
		{Name: "Sci-Fi", Description: "Science fiction content", SortOrder: 7},
		{Name: "Horror", Description: "Horror and thriller content", SortOrder: 8},
	}

	for _, cat := range categories {
		db.Create(&cat)
	}
	log.Printf("Seeded %d default categories", len(categories))
}

func seedAdmin(db *gorm.DB, cfg *config.Config) {
	var count int64
	db.Model(&models.User{}).Where("is_admin = ?", true).Count(&count)
	if count > 0 {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := models.User{
		Email:            cfg.AdminEmail,
		PasswordHash:     string(hashedPassword),
		Name:             "Admin",
		SubscriptionTier: "premium",
		IsAdmin:          true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("Admin user created: %s", cfg.AdminEmail)
}
