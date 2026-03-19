package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shyam/ott-backend/internal/database"
	"github.com/shyam/ott-backend/internal/models"
	"github.com/shyam/ott-backend/internal/utils"
)

type VideoHandler struct{}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{}
}

type CreateVideoRequest struct {
	Title              string `json:"title" binding:"required"`
	Description        string `json:"description"`
	YoutubeURL         string `json:"youtube_url" binding:"required"`
	CustomThumbnailURL string `json:"custom_thumbnail_url"`
	CategoryID         uint   `json:"category_id"`
	Duration           string `json:"duration"`
	ReleaseDate        string `json:"release_date"`
	IsPremium          bool   `json:"is_premium"`
	IsPublished        bool   `json:"is_published"`
}

type UpdateVideoRequest struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	YoutubeURL         string `json:"youtube_url"`
	CustomThumbnailURL string `json:"custom_thumbnail_url"`
	CategoryID         uint   `json:"category_id"`
	Duration           string `json:"duration"`
	ReleaseDate        string `json:"release_date"`
	IsPremium          *bool  `json:"is_premium"`
	IsPublished        *bool  `json:"is_published"`
}

// ListVideos returns published videos with optional filters
func (h *VideoHandler) ListVideos(c *gin.Context) {
	var videos []models.Video
	query := database.DB.Preload("Category").Where("is_published = ?", true)

	// Filter by category
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Search by title
	if search := c.Query("search"); search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	// Filter premium content based on user subscription
	user, _ := c.Get("user")
	u := user.(models.User)
	if u.SubscriptionTier != "premium" {
		if onlyFree := c.Query("only_free"); onlyFree == "true" {
			query = query.Where("is_premium = ?", false)
		}
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var total int64
	query.Model(&models.Video{}).Count(&total)

	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&videos)

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetVideo returns a single video by ID and increments view count
func (h *VideoHandler) GetVideo(c *gin.Context) {
	id := c.Param("id")

	var video models.Video
	if err := database.DB.Preload("Category").First(&video, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Increment view count
	database.DB.Model(&video).UpdateColumn("view_count", video.ViewCount+1)

	c.JSON(http.StatusOK, gin.H{"video": video})
}

// ContinueWatching returns videos the user has partially watched
func (h *VideoHandler) ContinueWatching(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var history []models.WatchHistory
	database.DB.Preload("Video").Preload("Video.Category").
		Where("user_id = ? AND watched_duration > 0 AND watched_duration < total_duration", userID).
		Order("last_watched_at DESC").
		Limit(20).
		Find(&history)

	c.JSON(http.StatusOK, gin.H{"continue_watching": history})
}

// Admin: Create video
func (h *VideoHandler) AdminCreateVideo(c *gin.Context) {
	var req CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	videoID, err := utils.ExtractYouTubeVideoID(req.YoutubeURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YouTube URL"})
		return
	}

	video := models.Video{
		Title:              req.Title,
		Description:        req.Description,
		YoutubeVideoID:     videoID,
		ThumbnailURL:       utils.GetThumbnailURL(videoID),
		CustomThumbnailURL: req.CustomThumbnailURL,
		CategoryID:         req.CategoryID,
		Duration:           req.Duration,
		ReleaseDate:        req.ReleaseDate,
		IsPremium:          req.IsPremium,
		IsPublished:        req.IsPublished,
	}

	if err := database.DB.Create(&video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video"})
		return
	}

	database.DB.Preload("Category").First(&video, video.ID)
	c.JSON(http.StatusCreated, gin.H{"video": video})
}

// Admin: Update video
func (h *VideoHandler) AdminUpdateVideo(c *gin.Context) {
	id := c.Param("id")

	var video models.Video
	if err := database.DB.First(&video, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	var req UpdateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.YoutubeURL != "" {
		videoID, err := utils.ExtractYouTubeVideoID(req.YoutubeURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YouTube URL"})
			return
		}
		updates["youtube_video_id"] = videoID
		updates["thumbnail_url"] = utils.GetThumbnailURL(videoID)
	}
	if req.CustomThumbnailURL != "" {
		updates["custom_thumbnail_url"] = req.CustomThumbnailURL
	}
	if req.CategoryID != 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Duration != "" {
		updates["duration"] = req.Duration
	}
	if req.ReleaseDate != "" {
		updates["release_date"] = req.ReleaseDate
	}
	if req.IsPremium != nil {
		updates["is_premium"] = *req.IsPremium
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}

	database.DB.Model(&video).Updates(updates)
	database.DB.Preload("Category").First(&video, video.ID)
	c.JSON(http.StatusOK, gin.H{"video": video})
}

// Admin: Delete video
func (h *VideoHandler) AdminDeleteVideo(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Video{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}

// Admin: List all videos (including unpublished)
func (h *VideoHandler) AdminListVideos(c *gin.Context) {
	var videos []models.Video
	query := database.DB.Preload("Category")

	// Search by title
	if search := c.Query("search"); search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	query.Order("created_at DESC").Find(&videos)
	c.JSON(http.StatusOK, gin.H{"videos": videos})
}
