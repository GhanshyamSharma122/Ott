package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shyam/ott-backend/internal/database"
	"github.com/shyam/ott-backend/internal/models"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetProfile returns the current user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}

	var user models.User
	database.DB.First(&user, userID)
	database.DB.Model(&user).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Admin: List all users
func (h *UserHandler) AdminListUsers(c *gin.Context) {
	var users []models.User
	query := database.DB

	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Order("created_at DESC").Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Admin: Update user (toggle premium, admin status, etc.)
func (h *UserHandler) AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		SubscriptionTier string `json:"subscription_tier"`
		IsAdmin          *bool  `json:"is_admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.SubscriptionTier != "" {
		updates["subscription_tier"] = req.SubscriptionTier
	}
	if req.IsAdmin != nil {
		updates["is_admin"] = *req.IsAdmin
	}

	database.DB.Model(&user).Updates(updates)
	database.DB.First(&user, id)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Admin: Get stats for dashboard
func (h *UserHandler) AdminStats(c *gin.Context) {
	var totalUsers, premiumUsers, totalVideos, totalCategories int64

	database.DB.Model(&models.User{}).Count(&totalUsers)
	database.DB.Model(&models.User{}).Where("subscription_tier = ?", "premium").Count(&premiumUsers)
	database.DB.Model(&models.Video{}).Count(&totalVideos)
	database.DB.Model(&models.Category{}).Count(&totalCategories)

	var totalViews int64
	database.DB.Model(&models.Video{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)

	c.JSON(http.StatusOK, gin.H{
		"total_users":      totalUsers,
		"premium_users":    premiumUsers,
		"total_videos":     totalVideos,
		"total_categories": totalCategories,
		"total_views":      totalViews,
	})
}
