package models

import (
	"time"
)

type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Email            string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash     string    `json:"-" gorm:"not null"`
	Name             string    `json:"name" gorm:"not null"`
	AvatarURL        string    `json:"avatar_url"`
	SubscriptionTier string    `json:"subscription_tier" gorm:"default:free"`
	IsAdmin          bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
