package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shyam/ott-backend/internal/database"
	"github.com/shyam/ott-backend/internal/models"
)

type WatchlistHandler struct{}

func NewWatchlistHandler() *WatchlistHandler {
	return &WatchlistHandler{}
}

// GetWatchlist returns user's watchlist
func (h *WatchlistHandler) GetWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var items []models.WatchlistItem
	database.DB.Preload("Video").Preload("Video.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items)

	c.JSON(http.StatusOK, gin.H{"watchlist": items})
}

// AddToWatchlist adds a video to the user's watchlist
func (h *WatchlistHandler) AddToWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		VideoID uint `json:"video_id" binding:"required"`
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

	// Check if already in watchlist
	var existing models.WatchlistItem
	if err := database.DB.Where("user_id = ? AND video_id = ?", userID, req.VideoID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Video already in watchlist"})
		return
	}

	item := models.WatchlistItem{
		UserID:  userID.(uint),
		VideoID: req.VideoID,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to watchlist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Added to watchlist", "item": item})
}

// RemoveFromWatchlist removes a video from the user's watchlist
func (h *WatchlistHandler) RemoveFromWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	videoID := c.Param("video_id")

	result := database.DB.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&models.WatchlistItem{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from watchlist"})
}

// CheckWatchlist checks if a video is in user's watchlist
func (h *WatchlistHandler) CheckWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	videoID := c.Param("video_id")

	var count int64
	database.DB.Model(&models.WatchlistItem{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{"in_watchlist": count > 0})
}
