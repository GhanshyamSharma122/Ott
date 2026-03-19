package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shyam/ott-backend/internal/database"
	"github.com/shyam/ott-backend/internal/models"
)

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

// Subscribe handles the dummy premium subscription (₹0)
func (h *PaymentHandler) Subscribe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, _ := c.Get("user")
	u := user.(models.User)

	if u.SubscriptionTier == "premium" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are already a premium subscriber"})
		return
	}

	// Create payment record
	payment := models.Payment{
		UserID:   userID.(uint),
		Amount:   0,
		Currency: "INR",
		Status:   "success",
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
		return
	}

	// Upgrade user to premium
	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("subscription_tier", "premium")

	// Fetch updated user
	var updatedUser models.User
	database.DB.First(&updatedUser, userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully upgraded to Premium!",
		"payment": payment,
		"user":    updatedUser,
	})
}

// GetSubscriptionStatus returns the current subscription status
func (h *PaymentHandler) GetSubscriptionStatus(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(models.User)

	var payments []models.Payment
	database.DB.Where("user_id = ?", u.ID).Order("created_at DESC").Find(&payments)

	c.JSON(http.StatusOK, gin.H{
		"subscription_tier": u.SubscriptionTier,
		"is_premium":        u.SubscriptionTier == "premium",
		"payments":          payments,
	})
}
