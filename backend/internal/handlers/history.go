package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shyam/ott-backend/internal/database"
	"github.com/shyam/ott-backend/internal/models"
)

type HistoryHandler struct{}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{}
}

// GetHistory returns user's watch history
func (h *HistoryHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var history []models.WatchHistory
	database.DB.Preload("Video").Preload("Video.Category").
		Where("user_id = ?", userID).
		Order("last_watched_at DESC").
		Limit(50).
		Find(&history)

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// UpdateProgress updates or creates a watch progress record
func (h *HistoryHandler) UpdateProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		VideoID         uint    `json:"video_id" binding:"required"`
		WatchedDuration float64 `json:"watched_duration" binding:"required"`
		TotalDuration   float64 `json:"total_duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if video exists
	var video models.Video
	if err := database.DB.First(&video, req.VideoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	var history models.WatchHistory
	result := database.DB.Where("user_id = ? AND video_id = ?", userID, req.VideoID).First(&history)

	now := time.Now()

	if result.Error != nil {
		// Create new history entry
		history = models.WatchHistory{
			UserID:          userID.(uint),
			VideoID:         req.VideoID,
			WatchedDuration: req.WatchedDuration,
			TotalDuration:   req.TotalDuration,
			LastWatchedAt:   now,
		}
		database.DB.Create(&history)
	} else {
		// Update existing entry
		database.DB.Model(&history).Updates(map[string]interface{}{
			"watched_duration": req.WatchedDuration,
			"total_duration":   req.TotalDuration,
			"last_watched_at":  now,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress updated", "history": history})
}
